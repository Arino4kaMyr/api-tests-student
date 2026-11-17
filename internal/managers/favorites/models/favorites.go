package models

import "api-tests-template/internal/managers/advertisements/models"

// AddFavoriteRequest тело запроса на добавление объявления в избранное
type AddFavoriteRequest struct {
	AdvertisementID string `json:"advertisement_id"`
}

// FavoriteItem запись избранного пользователя
type FavoriteItem struct {
	FavoriteID    string               `json:"favorite_id"`
	AddedAt       string               `json:"added_at"`
	Advertisement models.Advertisement `json:"advertisement"`
}

// FavoritesList страница списка избранного
type FavoritesList struct {
	Items      []FavoriteItem `json:"items"`
	Page       int            `json:"page"`
	PageSize   int            `json:"page_size"`
	Total      int            `json:"total"`
	TotalPages int            `json:"total_pages"`
}
