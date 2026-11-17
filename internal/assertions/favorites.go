package assertions

import "api-tests-template/internal/managers/favorites/models"

// CountFavoritesByAdvertisementID считает записи с указанным ID объявления.
func CountFavoritesByAdvertisementID(items []models.Favorite, advertisementID string) int {
	count := 0
	for _, item := range items {
		if item.Advertisement.ID == advertisementID {
			count++
		}
	}
	return count
}
