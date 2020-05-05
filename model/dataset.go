package model

type Dataset struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

// GetLabel will return the dataset's name. If the dataset does not have a title, it will instead return the ID
func (d Dataset) GetLabel() string {
	if d.Title == "" {
		return d.ID
	}
	return d.Title
}
