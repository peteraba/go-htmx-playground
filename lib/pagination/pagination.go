package pagination

import "strings"

type Pagination struct {
	Prev        int
	Next        int
	CurrentPage int
	Path        string
	Beginning   []int
	PreActive   []int
	PostActive  []int
	End         []int
}

func New(currentPage, pageSize, count int, path string) Pagination {
	if strings.Contains(path, "?") {
		path += "&page="
	} else {
		path += "?page="
	}

	maxPage := count / pageSize
	if count%pageSize > 0 {
		maxPage++
	}
	if maxPage == 0 {
		maxPage++
	}
	if currentPage > maxPage {
		currentPage = maxPage
	}

	return generate(maxPage, currentPage, path)
}

func generate(maxPage, currentPage int, path string) Pagination {
	var start, preActive, postActive, end []int

	prev := currentPage
	if currentPage > 1 {
		prev = currentPage - 1
	}

	next := currentPage
	if currentPage < maxPage {
		next = currentPage + 1
	}

	if currentPage >= 5 {
		start = []int{1, 2}
	} else if currentPage >= 4 {
		start = []int{1}
	}

	if currentPage >= 3 {
		preActive = []int{currentPage - 2, currentPage - 1}
	} else if currentPage >= 2 {
		preActive = []int{currentPage - 1}
	}

	if currentPage <= maxPage-2 {
		postActive = []int{currentPage + 1, currentPage + 2}
	} else if currentPage <= maxPage-1 {
		postActive = []int{currentPage + 1}
	}

	if currentPage <= maxPage-4 {
		end = []int{maxPage - 1, maxPage}
	} else if currentPage <= maxPage-3 {
		end = []int{maxPage}
	}

	return Pagination{
		Prev:        prev,
		Next:        next,
		CurrentPage: currentPage,
		Path:        path,
		Beginning:   start,
		PreActive:   preActive,
		PostActive:  postActive,
		End:         end,
	}
}
