package dto

type Search struct {
	SearchQuery string `json:"searchQuery" binding:"required"`
	Email 	 string `json:"email" binding:"required"`
}