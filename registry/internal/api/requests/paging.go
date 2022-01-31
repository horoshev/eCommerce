package requests

import (
	"github.com/go-chi/chi"
	"net/http"
	"strconv"
)

const DefaultPageNumber = 0
const DefaultPageSize = 10

type PageRequest struct {
	Page int `json:"page"`
	Size int `json:"size"`
}

// ParsePageRequest extracting paging information from request query
func ParsePageRequest(r *http.Request) *PageRequest {
	request := new(PageRequest)
	request.Page = intFromQueryOrDefault(r, "page", DefaultPageNumber)
	request.Size = intFromQueryOrDefault(r, "size", DefaultPageSize)

	return request
}

// intFromQueryOrDefault returns integer value of the named query parameter when it is presented.
// Returns default value otherwise.
func intFromQueryOrDefault(r *http.Request, name string, defaultValue int) int {
	if p := chi.URLParam(r, name); len(p) > 0 {
		if n, err := strconv.Atoi(p); err == nil && n > 0 {
			return n
		}
	}

	return defaultValue
}
