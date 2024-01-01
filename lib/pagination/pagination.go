package pagination

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/samber/lo"
)

type Pagination struct {
	Prev           string
	Next           string
	CurrentPage    string
	Path           string
	Target         string
	Start          []string
	PreActive      []string
	PostActive     []string
	End            []string
	Last           string
	IsPrevDisabled bool
	IsNextDisabled bool
	Query          map[string]interface{}
	count          int
}

func (p Pagination) SelfLink() string {
	return fmt.Sprintf("%s%s", p.Path, p.CurrentPage)
}

func (p Pagination) FirstLink() string {
	if p.CurrentPage == "1" {
		return ""
	}

	return fmt.Sprintf("%s1", p.Path)
}

func (p Pagination) PrevLink() string {
	if p.IsPrevDisabled {
		return ""
	}

	return fmt.Sprintf("%s%s", p.Path, p.Prev)
}

func (p Pagination) NextLink() string {
	if p.IsNextDisabled {
		return ""
	}

	return fmt.Sprintf("%s%s", p.Path, p.Next)
}

func (p Pagination) LastLink() string {
	if p.Last == "1" {
		return ""
	}

	return fmt.Sprintf("%s%s", p.Path, p.Last)
}

func (p Pagination) Count() int {
	return p.count
}

func New(currentPage, pageSize, count int, path string, params map[string]string, target string) Pagination {
	s := lo.MapToSlice(params, func(k, v string) string {
		return fmt.Sprintf("%s=%s", k, v)
	})
	s = append(s, "page=")

	path = fmt.Sprintf("%s?%s", path, strings.Join(s, "&"))

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

	return generate(maxPage, currentPage, path, target, params, count)
}

func generate(maxPage, currentPage int, path, target string, params map[string]string, count int) Pagination {
	var start, preActive, postActive, end []string

	prevPage := getPrev(currentPage)
	nextPage := getNext(currentPage, maxPage)

	start = getStart(currentPage)
	preActive = getPreActive(currentPage)
	postActive = getPostActive(currentPage, maxPage)
	end = getEnd(currentPage, maxPage)

	query := make(map[string]interface{}, len(params))
	for k, v := range params {
		query[k] = v
	}

	return Pagination{
		Prev:           strconv.Itoa(prevPage),
		Next:           strconv.Itoa(nextPage),
		CurrentPage:    strconv.Itoa(currentPage),
		Path:           path,
		Target:         target,
		Start:          start,
		PreActive:      preActive,
		PostActive:     postActive,
		End:            end,
		Last:           strconv.Itoa(maxPage),
		IsPrevDisabled: currentPage == 1,
		IsNextDisabled: currentPage == maxPage,
		count:          count,
		Query:          query,
	}
}

func toIntStringSlice(numbers ...int) []string {
	result := make([]string, 0, len(numbers))

	for _, num := range numbers {
		result = append(result, strconv.Itoa(num))
	}

	return result
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

func getStart(currentPage int) []string {
	if currentPage >= 5 {
		return []string{"1", "2"}
	} else if currentPage >= 4 {
		return []string{"1"}
	}

	return nil
}

func getPreActive(currentPage int) []string {
	if currentPage >= 3 {
		return toIntStringSlice(currentPage-2, currentPage-1)
	} else if currentPage >= 2 {
		return toIntStringSlice(currentPage - 1)
	}

	return nil
}

func getPostActive(currentPage, maxPage int) []string {
	if currentPage <= maxPage-2 {
		return toIntStringSlice(currentPage+1, currentPage+2)
	} else if currentPage <= maxPage-1 {
		return toIntStringSlice(currentPage + 1)
	}

	return nil
}

func getEnd(currentPage, maxPage int) []string {
	if currentPage <= maxPage-4 {
		return toIntStringSlice(maxPage-1, maxPage)
	} else if currentPage <= maxPage-3 {
		return toIntStringSlice(maxPage)
	}

	return nil
}
