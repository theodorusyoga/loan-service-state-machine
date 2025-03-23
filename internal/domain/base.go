package domain

// PaginatedResponse is a generic structure for paginated data
type PaginatedResponse struct {
	Data       any            `json:"data"`
	Pagination PaginationInfo `json:"pagination"`
}

// PaginationInfo contains metadata about the pagination
type PaginationInfo struct {
	CurrentPage int   `json:"current_page"`
	PageSize    int   `json:"page_size"`
	TotalItems  int64 `json:"total_items"`
	TotalPages  int   `json:"total_pages"`
}
