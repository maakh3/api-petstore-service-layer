package models

type Pet struct {
	Id        int      `json:"id"`
	Name      string   `json:"name"`
	Status    string   `json:"status"`
	Category  Category `json:"category"`
	Tags      []Tag    `json:"tags"`
	PhotoUrls []string `json:"photoUrls"`
}

type Category struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type Tag struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}
