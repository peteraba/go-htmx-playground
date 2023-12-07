package pagination

import "strings"

type Pagination struct {
	Prev        int
	Next        int
	CurrentPage int
	Path        string
	Target      string
	Beginning   []int
	PreActive   []int
	PostActive  []int
	End         []int
}

func New(currentPage, pageSize, count int, path, target string) Pagination {
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

	return generate(maxPage, currentPage, path, target)
}

func generate(maxPage, currentPage int, path, target string) Pagination {
	var start, preActive, postActive, end []int

	prev := getPrev(currentPage)
	next := getNext(currentPage, maxPage)

	start = getStart(currentPage)
	preActive = getPreActive(currentPage)
	postActive = getPostActive(currentPage, maxPage)
	end = getEnd(currentPage, maxPage)

	return Pagination{
		Prev:        prev,
		Next:        next,
		CurrentPage: currentPage,
		Path:        path,
		Target:      target,
		Beginning:   start,
		PreActive:   preActive,
		PostActive:  postActive,
		End:         end,
	}
}

func getPrev(currentPage int) int {
	if currentPage > 1 {
		return currentPage - 1
	}

	return currentPage
}

func getNext(currentPage, maxPage int) int {
	if currentPage < maxPage {
		return currentPage + 1
	}

	return currentPage
}

func getStart(currentPage int) []int {
	if currentPage >= 5 {
		return []int{1, 2}
	} else if currentPage >= 4 {
		return []int{1}
	}

	return nil
}

func getPreActive(currentPage int) []int {
	if currentPage >= 3 {
		return []int{currentPage - 2, currentPage - 1}
	} else if currentPage >= 2 {
		return []int{currentPage - 1}
	}

	return nil
}

func getPostActive(currentPage, maxPage int) []int {
	if currentPage <= maxPage-2 {
		return []int{currentPage + 1, currentPage + 2}
	} else if currentPage <= maxPage-1 {
		return []int{currentPage + 1}
	}

	return nil
}

func getEnd(currentPage, maxPage int) []int {
	if currentPage <= maxPage-4 {
		return []int{maxPage - 1, maxPage}
	} else if currentPage <= maxPage-3 {
		return []int{maxPage}
	}

	return nil
}
