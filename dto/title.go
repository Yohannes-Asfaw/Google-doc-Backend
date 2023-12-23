package dto

type Title struct {
	ID    string `json:"id" bson:"_id,omitempty"`
	Title string `json:"updatedTitle" binding:"required"`
}
