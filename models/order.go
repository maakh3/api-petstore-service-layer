package models

type Order struct {
	ID       int64  `json:"id"`
	PetID    int64  `json:"petId"`
	Quantity int    `json:"quantity"`
	ShipDate string `json:"shipDate"`
	Status   string `json:"status"` // Order Status: placed, approved, delivered
	Complete bool   `json:"complete"`
}
