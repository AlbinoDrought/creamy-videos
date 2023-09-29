package tmpl

import "fmt"

type PublicURLGenerator func(relativeURL string) string

type AppState struct {
	ReadOnly      bool
	SortDirection string
	Sortable      bool
	SearchText    string

	XSRFToken func() string

	PUG PublicURLGenerator
}

type Paging struct {
	URL         func(p int) string
	CurrentPage int
	Pages       int
}

type VideoFormState struct {
	Error       string
	Title       string
	Tags        string
	Description string
}

type paginationPage struct {
	URL      string
	Active   bool
	Disabled bool
	Page     string
}

func nextPageLink(p Paging) string {
	if p.CurrentPage+1 > p.Pages {
		return ""
	}
	return p.URL(p.CurrentPage + 1)
}

func genPages(p Paging) []paginationPage {
	const pageRange = 3

	pages := []paginationPage{}

	start := p.CurrentPage - pageRange
	if start < 1 {
		start = 1
	}
	end := p.CurrentPage + pageRange // inclusive
	if end > p.Pages {
		end = p.Pages
	}

	// show first page if not covered by page range
	if start > 1 {
		pages = append(pages, paginationPage{
			Page: "1",
			URL:  p.URL(1),
		})
	}

	// show ... if there is a gap between first page and page range
	if start > 2 {
		pages = append(pages, paginationPage{
			Page:     "...",
			Disabled: true,
		})
	}

	// gen page range
	for i := start; i <= end; i++ {
		pages = append(pages, paginationPage{
			Page:   fmt.Sprintf("%v", i),
			Active: p.CurrentPage == i,
			URL:    p.URL(i),
		})
	}

	// show ... if there is a gap between page range and last page
	if end < p.Pages-1 {
		pages = append(pages, paginationPage{
			Page:     "...",
			Disabled: true,
		})
	}

	// show last page if not covered by page range
	if end < p.Pages {
		pages = append(pages, paginationPage{
			Page: fmt.Sprintf("%v", p.Pages),
			URL:  p.URL(p.Pages),
		})
	}

	return pages
}
