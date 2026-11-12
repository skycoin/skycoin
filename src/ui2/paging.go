// Package ui2 src/ui2/paging.go c3-ui-wallet
package ui2

import (
	"strconv"
	"strings"
)

// ListPageSize is how many rows a long list renders at once.
const ListPageSize = 100

// ListQuery is the filter text and page a long list was asked for. The pages
// send it back with every request that redraws the list.
type ListQuery struct {
	Q    string
	Hide bool
	Page int
}

// ParseListQuery reads q, hide and page through get, which is usually a
// request's FormValue.
func ParseListQuery(get func(string) string) ListQuery {
	q := ListQuery{Q: strings.TrimSpace(get("q")), Hide: get("hide") != ""}
	q.Page, _ = strconv.Atoi(get("page")) //nolint:errcheck
	return q
}

// Match reports whether the filter text is in any of fields, ignoring case.
func (q ListQuery) Match(fields ...string) bool {
	if q.Q == "" {
		return true
	}
	needle := strings.ToLower(q.Q)
	for _, f := range fields {
		if strings.Contains(strings.ToLower(f), needle) {
			return true
		}
	}
	return false
}

// Pager is one page of a filtered list.
type Pager struct {
	ListQuery
	Total int
	Size  int
}

// Paginate keeps the rows of the page q asks for, moving an out of range page
// to the nearest one.
func Paginate[T any](rows []T, q ListQuery, size int) ([]T, Pager) {
	p := Pager{ListQuery: q, Total: len(rows), Size: size}
	p.Page = min(max(p.Page, 1), p.Pages())
	return rows[max(p.From()-1, 0):p.To()], p
}

// Pages is the number of pages, at least one.
func (p Pager) Pages() int {
	if p.Size <= 0 || p.Total == 0 {
		return 1
	}
	return (p.Total + p.Size - 1) / p.Size
}

// From is the 1 based index of the page's first row, 0 for an empty list.
func (p Pager) From() int {
	if p.Total == 0 {
		return 0
	}
	return (p.Page-1)*p.Size + 1
}

// To is the 1 based index of the page's last row.
func (p Pager) To() int {
	if p.Size <= 0 {
		return p.Total
	}
	return min(p.Page*p.Size, p.Total)
}

func (p Pager) pageVals(page int) string {
	return `{"page":"` + strconv.Itoa(page) + `"}`
}
