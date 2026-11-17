package models

// CreateAdvertisementRequest параметры создания объявления (multipart/form-data)
type CreateAdvertisementRequest struct {
	Title       string
	Description string
	Price       float64
	Quantity    int
	PhotosCount int // количество прикладываемых фотографий
}
