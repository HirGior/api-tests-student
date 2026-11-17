package models

type Advertisement struct {
	ID          string  `json:"id"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Quantity    int     `json:"quantity"`
	UserID      string  `json:"user_id"`
	Photos      []Photo `json:"photos"`
}

type AdvertisementsList struct {
	Items      []Advertisement `json:"items"`
	Total      int             `json:"total"`
	HasNext    bool            `json:"has_next"`
	NextCursor string          `json:"next_cursor"`
}

// Photo описывает фотографию, сохранённую на сервере.
type Photo struct {
	ID        string `json:"id"`
	URL       string `json:"url"`
	SortOrder int    `json:"sort_order"`
}
