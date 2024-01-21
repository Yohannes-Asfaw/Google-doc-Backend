// package main

// import (
// 	"fmt"
// 	"net/http"
// 	"net/http/httputil"
// 	"net/url"
// 	"os"

// 	"github.com/gorilla/websocket"
// )

// type Server interface {
// 	Address() string
// 	IsAlive() bool
// 	Serve(rw http.ResponseWriter, req *http.Request)
// }

// type simpleServer struct {
// 	addr  string
// 	proxy *httputil.ReverseProxy
// }

// func newSimpleServer(addr string) *simpleServer {
// 	serverUrl, err := url.Parse(addr)
// 	handleErr(err)

// 	return &simpleServer{
// 		addr:  addr,
// 		proxy: httputil.NewSingleHostReverseProxy(serverUrl),
// 	}
// }

// type LoadBalancer struct {
// 	port                        string
// 	roundRobinCountForHttp      int
// 	roundRobinCountForWebSocket int
// 	servers                     []Server
// 	//  documentId to connection mapping
// 	documentWebSockets map[string]*websocket.Conn
// 	// connection to server address mapping
// 	Connections map[*websocket.Conn]string
// }

// func NewLoadBalancer(port string, servers []Server) *LoadBalancer {
// 	return &LoadBalancer{
// 		port:                        port,
// 		roundRobinCountForHttp:      0,
// 		roundRobinCountForWebSocket: 0,
// 		servers:                     servers,
// 		documentWebSockets:          make(map[string]*websocket.Conn),
// 		Connections:                 make(map[*websocket.Conn]string),
// 	}
// }

// func handleErr(err error) {
// 	if err != nil {
// 		fmt.Printf("Error: %v\n", err)
// 		os.Exit(1)
// 	}
// }

// func (s *simpleServer) Address() string { return s.addr }

// func (s *simpleServer) IsAlive() bool {
// 	healthURL := s.addr + "/health"
// 	response, err := http.Get(healthURL)
// 	if err != nil {
// 		fmt.Printf("Error checking server %s health: %v\n", s.addr, err)
// 		return false
// 	}
// 	defer response.Body.Close()

// 	return response.StatusCode == http.StatusOK
// }


// func (s *simpleServer) Serve(rw http.ResponseWriter, req *http.Request) {
// 	s.proxy.ServeHTTP(rw, req)
// }

// func (lb *LoadBalancer) getNextAvailableServer(isSocket bool) Server {
// 	loop := 0
// 	if isSocket {
// 		server := lb.servers[lb.roundRobinCountForWebSocket%len(lb.servers)]
// 		for !server.IsAlive() {
// 			lb.roundRobinCountForWebSocket++
// 			loop++
// 			server = lb.servers[lb.roundRobinCountForWebSocket%len(lb.servers)]
// 			if loop > len(lb.servers) {
// 				fmt.Print("No server is alive\n")
// 				loop = 0
// 				// os.Exit(1)
// 			}
// 		}
// 		// fmt.Print("Round Robin socket Count is : ", lb.roundRobinCountForWebSocket, ", Selected Server: ", server.Address(), "\n")
// 		lb.roundRobinCountForWebSocket++
// 		return server
// 	} else {
// 		server := lb.servers[lb.roundRobinCountForHttp%len(lb.servers)]
// 		for !server.IsAlive() {
// 			lb.roundRobinCountForHttp++
// 			loop++
// 			server = lb.servers[lb.roundRobinCountForHttp%len(lb.servers)]
// 			if loop > len(lb.servers) {
// 				fmt.Print("No server is alive\n")
// 				loop = 0
// 				// os.Exit(1)
// 			}
// 		}
// 		// fmt.Print("Round Robin http Count is : ", lb.roundRobinCountForHttp, ", Selected Server: ", server.Address(), "\n")
// 		lb.roundRobinCountForHttp++
// 		return server
// 	}

// }

// func (lb *LoadBalancer) getServerWithExistingConnection(documentID string) Server {

// 	// Check if there is an existing WebSocket connection for the document
// 	if conn, ok := lb.documentWebSockets[documentID]; ok {
// 		// Check if the connection has a corresponding server address
// 		if serverAddr, ok := lb.Connections[conn]; ok {
// 			// Find the server object with the matching address
// 			for _, s := range lb.servers {
// 				if s.Address() == serverAddr && s.IsAlive() {
// 					return s
// 				}
// 			}
// 		}
// 	}

// 	// If no existing connection is found, create a new one
// 	server := lb.getNextAvailableServer(true)
// 	lb.documentWebSockets[documentID] = &websocket.Conn{}
// 	lb.Connections[lb.documentWebSockets[documentID]] = server.Address()
// 	return server
// }

// func (lb *LoadBalancer) serveProxy(rw http.ResponseWriter, req *http.Request) {
// 	documentID := req.URL.Query().Get("document_id")
// 	if documentID == "" {
// 		targetServer := lb.getNextAvailableServer(false)
// 		targetServer.Serve(rw, req)
// 	} else {
// 		targetServer := lb.getServerWithExistingConnection(documentID)
// 		targetServer.Serve(rw, req)
// 	}
// }

// func main() {
// 	servers := []Server{
// 		newSimpleServer("http://127.0.0.1:8080"),
// 		newSimpleServer("http://127.0.0.1:8081"),
// 		newSimpleServer("http://127.0.0.1:8082"),
// 	}

// 	lb := NewLoadBalancer("7000", servers)
// 	handleRedirect := func(rw http.ResponseWriter, req *http.Request) {
// 		lb.serveProxy(rw, req)
// 	}

// 	http.HandleFunc("/", handleRedirect)
// 	http.ListenAndServe("127.0.0.1:"+lb.port, nil)
// }
