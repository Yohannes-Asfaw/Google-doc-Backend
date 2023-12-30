package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/khallihub/godoc/controller"
	"github.com/khallihub/godoc/dto"
	"github.com/khallihub/godoc/middlewares"
	"github.com/khallihub/godoc/service"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"github.com/joho/godotenv"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type DocumentWebSocket struct {
	Connections map[*websocket.Conn]bool
	Mutex       sync.Mutex
}

var documentWebSockets = make(map[string]*DocumentWebSocket)
var documentCache sync.Map

func main() {

	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file:", err)
		return
	}

	dbUrl := os.Getenv("DATABASE_URL")

	if dbUrl == "" {
		fmt.Println("DATABASE_URL not found in .env file")
		return
	}

	// MongoDB connection setup
	mongoClient, err := mongo.NewClient(options.Client().ApplyURI(dbUrl))
	if err != nil {
		panic(err)
	}
	err = mongoClient.Connect(context.Background())
	if err != nil {
		panic(err)
	}
	defer mongoClient.Disconnect(context.Background())

	server := gin.New()
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Authorization", "Content-Type"}

	server.Use(cors.New(config))

	server.Use(gin.Recovery(), gin.Logger())

	signupService := service.NewSignupService(mongoClient, "godoc", "users")
	signupController := controller.NewSignupController(signupService)

	loginService := service.NewLoginService(mongoClient, "godoc", "users")
	jwtService := service.NewJWTService()
	loginController := controller.NewLoginController(loginService, jwtService)

	// Route for chaecking the health of the server
	server.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"status": "UP",
		})
	})

	// Routes for handling user authentication
	authRoutes := server.Group("/auth")
	{
		// Signup Endpoint: User creation
		authRoutes.POST("/signup", func(ctx *gin.Context) {
			message := signupController.Signup(ctx)
			if message != "" {
				ctx.JSON(http.StatusOK, gin.H{
					"message": message,
				})
			} else {
				ctx.JSON(http.StatusBadRequest, nil)
			}
		})

		// Login Endpoint: Authentication + Token creation
		authRoutes.POST("/login", func(ctx *gin.Context) {
			token := loginController.Login(ctx)
			if token != "" {
				ctx.JSON(http.StatusOK, gin.H{
					"token": token,
				})
			} else {
				ctx.JSON(http.StatusUnauthorized, gin.H{
					"message": "Invalid credentials",
				})
			}
		})
	}

	documentService := service.NewDocumentService(mongoClient, "godoc", "documents")
	documentController := controller.NewDocumentController(documentService)

	// Route for handling document operations
	documentRoutes := server.Group(("/documents"))
	documentRoutes.Use(middlewares.AuthorizeJWT())
	{
		documentRoutes.GET("/handler", func(ctx *gin.Context) {
			documentID := ctx.Query("document_id")
			handleWebSocket(ctx, documentID, documentController)
		})

		// Route for getting all documents
		documentRoutes.POST("/getall", func(ctx *gin.Context) {
			// Fetching documents from MongoDB and responding with JSON
			documents, err := documentController.GetAllDocuments(ctx)
			if err != nil {
				return
			}
			ctx.JSON(http.StatusOK, gin.H{
				"documents": documents,
			})
		})

		// Route for seraching documents by title
		documentRoutes.POST("/search", func(ctx *gin.Context) {
			documents, err := documentController.SearchDocuments(ctx)
			if err != nil {
				fmt.Println("Error searching documents:", err)
				return
			}
			ctx.JSON(http.StatusOK, gin.H{
				"documents": documents,
			})
		})

		// Route for creating a new document
		documentRoutes.POST("/createnew", func(ctx *gin.Context) {
			// Creating a new document and storing it in MongoDB
			documentController.CreateNewDocument(ctx)
		})

		// Route for getting a specific document
		documentRoutes.POST("/getone/:id", func(ctx *gin.Context) {

			document, err := initializeDocumentCache(ctx, documentController)

			if err != nil {
				fmt.Println("Error getting document:", err)
				return
			}
			ctx.JSON(http.StatusOK, document)
		})

		documentRoutes.POST("/updatetitle", func(ctx *gin.Context) {
			// Updating the title of a document
			title, documentID := documentController.UpdateTitle(ctx)
			updateDocumentTitleCacheAttribute(documentID, title)
		})

		documentRoutes.POST("/updatecollaborators", func(ctx *gin.Context) {
			// Adding a collaborator to a document
			// updating the database
			document := documentController.UpdateCollaborators(ctx)
			// updating the cache
			access := new(dto.Access)
			access.ID = document.ID
			access.ReadAccess = document.ReadAccess
			access.WriteAccess = document.WriteAccess
			updateDocumentCacheAttribute(document.ID, documentController, *access)
		})

		// Route for deleting a document
		documentRoutes.DELETE("/delete/:id", func(ctx *gin.Context) {
			// Deleting a document from MongoDB
			documentController.DeleteDocument(ctx)
		})
	}

	// Start the periodic cache update
	updateDatabaseWithCache(documentController)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	server.Run("127.0.0.1:" + port)
}

func handleWebSocket(ctx *gin.Context, documentID string, documentController controller.DocumentController) {
	fmt.Println("Handling WebSocket connection for document:", documentID)
	fmt.Println("Connection handled by server running on port:", os.Getenv("PORT"))
	
	upgrader.CheckOrigin = func(r *http.Request) bool {
		// Allow any origin (not recommended for production, consider a more restrictive check)
		return true
	}
	
	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		log.Println("Error upgrading to WebSocket:", err)
		return
	}
	defer conn.Close()
	
	// Get or create a WebSocket instance for the documentID
	documentWebSocket, ok := documentWebSockets[documentID]
	if !ok {
		documentWebSocket = &DocumentWebSocket{Connections: make(map[*websocket.Conn]bool)}
		documentWebSockets[documentID] = documentWebSocket
	}
	fmt.Println("Number of active connections:", len(documentWebSockets[documentID].Connections)+1)

	// Add the new connection to the WebSocket instance
	documentWebSocket.Mutex.Lock()
	documentWebSocket.Connections[conn] = true
	documentWebSocket.Mutex.Unlock()

	// Create a channel to signal when a client disconnects
	disconnectChannel := make(chan *websocket.Conn, 1)

	// Save the source connection
	sourceConnection := conn

	// Start a goroutine to handle disconnection cleanup
	go func() {
		select {
		case <-disconnectChannel:
			// Remove the connection from the WebSocket instance when the client disconnects
			documentWebSocket.Mutex.Lock()
			delete(documentWebSocket.Connections, sourceConnection)
			if len(documentWebSocket.Connections) == 0 {
				fmt.Println("No more connections. Cleaning up resources for document:", documentID)
				documentCache.Delete(documentID)
				delete(documentWebSockets, documentID)
			}
			documentWebSocket.Mutex.Unlock()
		}
	}()

	// Save the source connection
	// sourceConnection := conn

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error reading message:", err)
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("Error reading message: %v", err)
			}
			break
		}

		var message dto.Message
		if err := json.Unmarshal(msg, &message); err != nil {
			log.Println("Error unmarshalling document:", err)
			continue
		}

		// Update the document cache
		if err := updateDocumentCache(documentID, documentController, message.Data); err != nil {
			log.Println("Error updating document cache:", err)
			continue
		}

		// Update the document in the database
		// if err := documentController.UpdateDocument(documentID, message.Data); err != nil {
		// 	log.Println("Error updating document in DB:", err)
		// 	continue
		// }

		// Broadcast the message to all connected clients for the document
		documentWebSocket.Mutex.Lock()
		for conn := range documentWebSocket.Connections {
			// Skip broadcasting to the source connection
			if conn == sourceConnection {
				continue
			}

			data, err := json.Marshal(message.Change)
			if err != nil {
				log.Println("Error marshalling message:", err)
				continue
			}

			err = conn.WriteMessage(websocket.TextMessage, data)
			if err != nil {
				log.Println("Error writing message:", err)
				conn.Close()
				disconnectChannel <- conn // Signal disconnection to the cleanup goroutine
				delete(documentWebSocket.Connections, conn)
			}
		}
		documentWebSocket.Mutex.Unlock()
	}
	close(disconnectChannel)
}

func initializeDocumentCache(ctx *gin.Context, documentController controller.DocumentController) (*dto.Document, error) {
	documentID := ctx.Param("id")
	var document *dto.Document

	// Check if the document is already in the cache
	if cachedDocument, ok := documentCache.Load(documentID); !ok {
		// Fetch the document from the database
		document, err := documentController.GetOneDocument(ctx)
		if err != nil {
			fmt.Println("Error getting document:", err)
			return nil, err
		}

		// Store the document in the cache
		documentCache.Store(documentID, document)
		ctx.JSON(http.StatusOK, document)
	} else {
		// If document is already in cache, retrieve it
		document = cachedDocument.(*dto.Document)
	}
	return document, nil
}

func updateDocumentCache(documentID string, documentController controller.DocumentController, newData dto.DocumentData) error {
	cachedDocument, ok := documentCache.Load(documentID)
	if !ok {
		return fmt.Errorf("document not found in cache")
	}

	document := cachedDocument.(*dto.Document)
	document.Data = newData

	// Update the document in the cache
	documentCache.Store(documentID, document)
	return nil
}

func updateDatabaseWithCache(documentController controller.DocumentController) {
	// Create a ticker that ticks every specified duration
	ticker := time.NewTicker(30 * time.Second)

	// Run a goroutine to perform the periodic update
	go func() {
		for {
			select {
			case <-ticker.C:
				// Perform the database update using the cache
				err := syncDatabaseWithCache(documentController)
				if err != nil {
					fmt.Println("Error updating database with cache:", err)
				}
			}
		}
	}()
}

func syncDatabaseWithCache(documentController controller.DocumentController) error {
	// Iterate over the cache and update the database with each entry
	documentCache.Range(func(key, value interface{}) bool {
		documentID := key.(string)
		cachedDocument := value.(*dto.Document)

		// Update the database with the cached document
		err := documentController.UpdateDocument(documentID, cachedDocument.Data)
		if err != nil {
			// Log or handle the error accordingly
			fmt.Printf("Error updating database for document %s: %v\n", documentID, err)
		}
		// documentCache = sync.Map{}
		return true
	})

	return nil
}

func updateDocumentCacheAttribute(documentID string, documentController controller.DocumentController, newData dto.Access) error {
	fmt.Print("Updating document cache attribute\n")
	cachedDocument, ok := documentCache.Load(documentID)
	if !ok {
		return fmt.Errorf("document not found in cache")
	}

	document := cachedDocument.(*dto.Document)
	document.ReadAccess = newData.ReadAccess
	document.WriteAccess = newData.WriteAccess

	// Update the document in the cache
	documentCache.Store(documentID, document)
	return nil
}

func updateDocumentTitleCacheAttribute(documentID string, newTitle string) error {
	fmt.Print("Updating document title cache attribute\n", documentID, newTitle)
	cachedDocument, ok := documentCache.Load(documentID)
	if !ok {
		return fmt.Errorf("document not found in cache")
	}

	document := cachedDocument.(*dto.Document)
	document.Title = newTitle

	// Update the document in the cache
	documentCache.Store(documentID, document)
	return nil
}