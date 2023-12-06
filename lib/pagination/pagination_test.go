//nolint:exhaustruct
package pagination_test

import (
	"reflect"
	"testing"

	"github.com/peteraba/go-htmx-playground/lib/pagination"
)

func TestNew(t *testing.T) {
	t.Parallel()

	type args struct {
		currentPage int
		pageSize    int
		count       int
		path        string
	}
	tests := []struct {
		name string
		args args
		want pagination.Pagination
	}{
		{
			name: "empty",
			args: args{
				currentPage: 1,
				pageSize:    10,
				count:       0,
				path:        "/hello",
			},
			want: pagination.Pagination{
				Next:        1,
				Prev:        1,
				CurrentPage: 1,
				Path:        "/hello?page=",
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
			want: pagination.Pagination{
				Next:        1,
				Prev:        1,
				CurrentPage: 1,
				Path:        "/hello?page=",
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
			want: pagination.Pagination{
				Next:        1,
				Prev:        1,
				CurrentPage: 1,
				Path:        "/hello?page=",
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
			want: pagination.Pagination{
				Next:        2,
				Prev:        1,
				CurrentPage: 1,
				Path:        "/hello?page=",
				PostActive:  []int{2},
			},
		},
	}
	for _, ttt := range tests {
		tt := ttt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := pagination.New(tt.args.currentPage, tt.args.pageSize, tt.args.count, tt.args.path); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}
