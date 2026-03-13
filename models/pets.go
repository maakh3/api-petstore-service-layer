package models

type Pet struct {
	Id     int    `json:"id"`
	Name   string `json:"name"`
	Status string `json:"status"`
}
