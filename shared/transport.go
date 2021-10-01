package shared

import (
	"net/http"
	"strconv"

	"github.com/unknwon/paginater"
)

type PaginatedRequest struct {
	Page          int
	PageSize      int
	Query         string
	SortLabel     string
	SortDirection string
}

func NewPaginatedRequest(r *http.Request) PaginatedRequest {
	q := r.URL.Query()

	page, _ := strconv.Atoi(q.Get("page"))
	pageSize, _ := strconv.Atoi(q.Get("page_size"))

	if page <= 0 {
		page = 1
	}

	if pageSize <= 1 {
		pageSize = 1
	}

	var pr PaginatedRequest

	pr.Page = page
	pr.PageSize = pageSize
	pr.Query = q.Get("query")
	pr.SortLabel = q.Get("sort")
	pr.SortDirection = q.Get("sort_dir")

	return pr
}

type PaginatedResponse struct {
	CurrentPage     int  `json:"current_page"`
	TotalPages      int  `json:"total_pages"`
	TotalCount      int  `json:"total_count"`
	PageSize        int  `json:"page_size"`
	PreviousPage    int  `json:"previous_page"`
	HasPreviousPage bool `json:"has_previous_page"`
	NextPage        int  `json:"next_page"`
	HasNextPage     bool `json:"has_next_page"`
}

func NewPaginatedResponse(total, size, current int) PaginatedResponse {
	pager := paginater.New(total, size, current, 5)

	var pr PaginatedResponse

	pr.CurrentPage = pager.Current()
	pr.PageSize = pager.PagingNum()
	pr.TotalCount = pager.Total()
	pr.TotalPages = pager.TotalPages()
	pr.PreviousPage = pager.Previous()
	pr.HasPreviousPage = pager.HasPrevious()
	pr.NextPage = pager.Next()
	pr.HasNextPage = pager.HasNext()

	return pr
}
