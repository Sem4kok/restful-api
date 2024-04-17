package util

// Album type for storing album's
type Album struct {
	ID     int    `json:"id"`
	Title  string `json:"title" binding:"required"`
	Artist string `json:"artist" binding:"required"`
	Price  string `json:"price"`
}
