package models

// CreateAdvertisementRequest содержит все поля создания объявления и пути к фотографиям.
type CreateAdvertisementRequest struct {
	Title       string
	Description string
	Price       float64
	Quantity    int
	Photos      []string
}
