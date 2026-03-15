package services

import "errors"

var (
	ErrPetNotFound   = errors.New("pet not found")
	ErrOrderNotFound = errors.New("order not found")
)
