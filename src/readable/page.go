package readable

// PageInfo represents the pagination info
type PageInfo struct {
	TotalPages  uint64 `json:"total_pages"`
	PageSize    uint64 `json:"page_size"`
	CurrentPage uint64 `json:"current_page"`
}
