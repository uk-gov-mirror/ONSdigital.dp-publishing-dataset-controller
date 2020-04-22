package model

type Dataset struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

// GetTitle will return the dataset's title. If the dataset does not have a title, it will instead return the ID
func (d Dataset) GetTitle() string {
	if d.Title == "" {
		return d.ID
	}
	return d.Title
}
