package utils

type PaginatedResponse[T any] struct {
	Page       int64 `json:"page"`
	PageSize   int64 `json:"page_size"`
	TotalItems int64 `json:"total_items"`
	Items      []T   `json:"items"`
}
