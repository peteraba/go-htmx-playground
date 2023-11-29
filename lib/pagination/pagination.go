package pagination

type Pagination struct {
	Prev        bool
	Next        bool
	CurrentPage int
	Path        string
	Start       []int
	PreActive   []int
	PostActive  []int
	End         []int
}

func New(currentPage, pageSize, count int, path string) Pagination {
	maxPage := (count / pageSize)
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
		Prev:        currentPage > 1,
		Next:        currentPage < maxPage,
		CurrentPage: currentPage,
		Path:        path,
		Start:       start,
		PreActive:   preActive,
		PostActive:  postActive,
		End:         end,
	}
}
