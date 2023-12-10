package dto

type DocumentData struct {
	Ops []map[string]interface{} `json:"ops" bson:"ops"`
}

type Document struct {
	ID     string        `json:"id" bson:"_id,omitempty"`
	Author string        `json:"author" binding:"required"`
	ReadAccess []string      `json:"readAccess" bson:"readAccess"`
	WriteAccess []string      `json:"writeAccess" bson:"writeAccess"`
	Title  string        `json:"title"`
	Data   DocumentData `json:"data" bson:"data"`
}

type Message struct {
	Data   DocumentData `json:"data" bson:"data"`
	Change map[string]interface{} 
}