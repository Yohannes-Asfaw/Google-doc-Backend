package dto

type Access struct {
	ID     string `json:"document_id" bson:"_id"`
	ReadAccess []string `json:"readAccess" bson:"readAccess"`
	WriteAccess []string `json:"writeAccess" bson:"writeAccess"`
}
