package pagination

import (
	"reflect"
	"testing"
)

func Test_generate(t *testing.T) {
	type args struct {
		maxPage     int
		currentPage int
		path        string
	}
	tests := []struct {
		name string
		args args
		want Pagination
	}{
		{
			name: "default",
			args: args{
				maxPage:     1,
				currentPage: 1,
				path:        "/hello",
			},
			want: Pagination{
				CurrentPage: 1,
				Path:        "/hello",
			},
		},
		{
			name: "1/2",
			args: args{
				maxPage:     2,
				currentPage: 1,
				path:        "/hello",
			},
			want: Pagination{
				Next:        true,
				CurrentPage: 1,
				Path:        "/hello",
				PostActive:  []int{2},
			},
		},
		{
			name: "1/3",
			args: args{
				maxPage:     3,
				currentPage: 1,
				path:        "/hello",
			},
			want: Pagination{
				Next:        true,
				CurrentPage: 1,
				Path:        "/hello",
				PostActive:  []int{2, 3},
			},
		},
		{
			name: "1/4",
			args: args{
				maxPage:     4,
				currentPage: 1,
				path:        "/hello",
			},
			want: Pagination{
				Next:        true,
				CurrentPage: 1,
				Path:        "/hello",
				PostActive:  []int{2, 3},
				End:         []int{4},
			},
		},
		{
			name: "1/5",
			args: args{
				maxPage:     5,
				currentPage: 1,
				path:        "/hello",
			},
			want: Pagination{
				Next:        true,
				CurrentPage: 1,
				Path:        "/hello",
				PostActive:  []int{2, 3},
				End:         []int{4, 5},
			},
		},
		{
			name: "2/5",
			args: args{
				maxPage:     5,
				currentPage: 2,
				path:        "/hello",
			},
			want: Pagination{
				Prev:        true,
				Next:        true,
				CurrentPage: 2,
				Path:        "/hello",
				PreActive:   []int{1},
				PostActive:  []int{3, 4},
				End:         []int{5},
			},
		},
		{
			name: "3/5",
			args: args{
				maxPage:     5,
				currentPage: 3,
				path:        "/hello",
			},
			want: Pagination{
				Prev:        true,
				Next:        true,
				CurrentPage: 3,
				Path:        "/hello",
				PreActive:   []int{1, 2},
				PostActive:  []int{4, 5},
			},
		},
		{
			name: "4/5",
			args: args{
				maxPage:     5,
				currentPage: 4,
				path:        "/hello",
			},
			want: Pagination{
				Prev:        true,
				Next:        true,
				CurrentPage: 4,
				Path:        "/hello",
				Start:       []int{1},
				PreActive:   []int{2, 3},
				PostActive:  []int{5},
			},
		},
		{
			name: "5/5",
			args: args{
				maxPage:     5,
				currentPage: 5,
				path:        "/hello",
			},
			want: Pagination{
				Prev:        true,
				Next:        false,
				CurrentPage: 5,
				Path:        "/hello",
				Start:       []int{1, 2},
				PreActive:   []int{3, 4},
			},
		},
		{
			name: "49/73",
			args: args{
				maxPage:     73,
				currentPage: 49,
				path:        "/hello",
			},
			want: Pagination{
				Prev:        true,
				Next:        true,
				CurrentPage: 49,
				Path:        "/hello",
				Start:       []int{1, 2},
				PreActive:   []int{47, 48},
				PostActive:  []int{50, 51},
				End:         []int{72, 73},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := generate(tt.args.maxPage, tt.args.currentPage, tt.args.path); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNew(t *testing.T) {
	type args struct {
		currentPage int
		pageSize    int
		count       int
		path        string
	}
	tests := []struct {
		name string
		args args
		want Pagination
	}{
		{
			name: "empty",
			args: args{
				currentPage: 1,
				pageSize:    10,
				count:       0,
				path:        "/hello",
			},
			want: Pagination{
				CurrentPage: 1,
				Path:        "/hello",
			},
		},
		{
			name: "1/1 min",
			args: args{
				currentPage: 1,
				pageSize:    10,
				count:       1,
				path:        "/hello",
			},
			want: Pagination{
				CurrentPage: 1,
				Path:        "/hello",
			},
		},
		{
			name: "1/1 max",
			args: args{
				currentPage: 1,
				pageSize:    10,
				count:       10,
				path:        "/hello",
			},
			want: Pagination{
				CurrentPage: 1,
				Path:        "/hello",
			},
		},
		{
			name: "1/2",
			args: args{
				currentPage: 1,
				pageSize:    10,
				count:       11,
				path:        "/hello",
			},
			want: Pagination{
				Next:        true,
				CurrentPage: 1,
				Path:        "/hello",
				PostActive:  []int{2},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(tt.args.currentPage, tt.args.pageSize, tt.args.count, tt.args.path); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}
