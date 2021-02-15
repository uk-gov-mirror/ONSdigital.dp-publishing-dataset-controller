package topics

type TopicsResult struct {
	Topics Topic `json:"topics"`
}

type Topic struct {
	Results []Result `json:"results"`
}

type Result struct {
	Description Description `json:"description"`
	URI         string      `json:"uri"`
	Type        string      `json:"type"`
}

type Description struct {
	Title string `json:"title"`
}
