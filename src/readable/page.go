// Package readable src/readable/page.go c2-coin-api
package readable

// PageInfo represents the pagination info
type PageInfo struct {
	TotalPages  uint64 `json:"total_pages"`
	PageSize    uint64 `json:"page_size"`
	CurrentPage uint64 `json:"current_page"`
}
