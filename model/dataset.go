package model

type Dataset struct {
	ID           string `json:"id"`
	CollectionID string `json:"collection_id"`
	Title        string `json:"title"`
	Description  string `json:"description"`
}
