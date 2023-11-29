package repository

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/peteraba/go-htmx-playground/pkg/films/model"
)

func TestFilmRepo_CountDirectors(t *testing.T) {
	type fields struct {
		films    []model.Film
		maxLimit int
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		{
			name: "empty",
			fields: fields{
				maxLimit: 10,
			},
			want: 0,
		},
		{
			name: "single film",
			fields: fields{
				films: []model.Film{
					{Title: "foo", Director: "bar"},
				},
				maxLimit: 10,
			},
			want: 1,
		},
		{
			name: "no directors with multiple films",
			fields: fields{
				films: []model.Film{
					{Title: "foo", Director: "bar"},
					{Title: "baz", Director: "quix"},
				},
				maxLimit: 10,
			},
			want: 2,
		},
		{
			name: "directors with multiple films",
			fields: fields{
				films: []model.Film{
					{Title: "foo", Director: "bar"},
					{Title: "baz", Director: "bar"},
				},
				maxLimit: 10,
			},
			want: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewFilmRepo(tt.fields.maxLimit, tt.fields.films...)
			if got := r.CountDirectors(); got != tt.want {
				t.Errorf("CountDirectors() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFilmRepo_CountFilms(t *testing.T) {
	type fields struct {
		films    []model.Film
		maxLimit int
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		{
			name: "empty",
			fields: fields{
				maxLimit: 10,
			},
			want: 0,
		},
		{
			name: "single film",
			fields: fields{
				films: []model.Film{
					{Title: "foo", Director: "bar"},
				},
				maxLimit: 10,
			},
			want: 1,
		},
		{
			name: "no directors with multiple films",
			fields: fields{
				films: []model.Film{
					{Title: "foo", Director: "bar"},
					{Title: "baz", Director: "quix"},
				},
				maxLimit: 10,
			},
			want: 2,
		},
		{
			name: "directors with multiple films",
			fields: fields{
				films: []model.Film{
					{Title: "foo", Director: "bar"},
					{Title: "baz", Director: "bar"},
				},
				maxLimit: 10,
			},
			want: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewFilmRepo(tt.fields.maxLimit, tt.fields.films...)
			if got := r.CountFilms(); got != tt.want {
				t.Errorf("CountFilms() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFilmRepo_Insert(t *testing.T) {
	type fields struct {
		films    []model.Film
		maxLimit int
	}
	type args struct {
		newFilms []model.Film
	}
	tests := []struct {
		name              string
		fields            fields
		args              args
		wantFilmCount     int
		wantDirectorCount int
	}{
		{
			name: "from empty",
			fields: fields{
				maxLimit: 10,
			},
			args: args{
				newFilms: []model.Film{
					{Title: "foo", Director: "bar"},
					{Title: "baz", Director: "bar"},
				},
			},
			wantFilmCount:     2,
			wantDirectorCount: 1,
		},
		{
			name: "from non-empty",
			fields: fields{
				maxLimit: 10,
				films: []model.Film{
					{Title: "quix", Director: "bar"},
				},
			},
			args: args{
				newFilms: []model.Film{
					{Title: "foo", Director: "bar"},
					{Title: "baz", Director: "bar"},
				},
			},
			wantFilmCount:     3,
			wantDirectorCount: 1,
		},
		{
			name: "duplicate titles get skipped",
			fields: fields{
				maxLimit: 10,
				films: []model.Film{
					{Title: "foo", Director: "quix"},
				},
			},
			args: args{
				newFilms: []model.Film{
					{Title: "foo", Director: "bar"},
					{Title: "baz", Director: "bar"},
				},
			},
			wantFilmCount:     2,
			wantDirectorCount: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewFilmRepo(tt.fields.maxLimit, tt.fields.films...)
			_ = r.Insert(tt.args.newFilms...)

			assert.Equal(t, tt.wantFilmCount, r.CountFilms())
			assert.Equal(t, tt.wantDirectorCount, r.CountDirectors())
		})
	}
}

func TestFilmRepo_ListDirectors(t *testing.T) {

	type fields struct {
		films    []model.Film
		maxLimit int
	}
	type args struct {
		offset int
		limit  int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []model.Director
	}{
		{
			name: "empty",
			fields: fields{
				maxLimit: 10,
			},
			args: args{
				offset: 0,
				limit:  10,
			},
		},
		{
			name: "single film",
			fields: fields{
				films: []model.Film{
					{Title: "foo", Director: "bar"},
				},
				maxLimit: 10,
			},
			args: args{
				offset: 0,
				limit:  10,
			},
			want: []model.Director{
				{Name: "bar", Titles: []string{"foo"}},
			},
		},
		{
			name: "no directors with multiple films",
			fields: fields{
				films: []model.Film{
					{Title: "foo", Director: "bar"},
					{Title: "baz", Director: "quix"},
				},
				maxLimit: 10,
			},
			args: args{
				offset: 0,
				limit:  10,
			},
			want: []model.Director{
				{Name: "bar", Titles: []string{"foo"}},
				{Name: "quix", Titles: []string{"baz"}},
			},
		},
		{
			name: "directors with multiple films",
			fields: fields{
				films: []model.Film{
					{Title: "foo", Director: "bar"},
					{Title: "baz", Director: "bar"},
				},
				maxLimit: 10,
			},
			args: args{
				offset: 0,
				limit:  10,
			},
			want: []model.Director{
				{Name: "bar", Titles: []string{"foo", "baz"}},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewFilmRepo(tt.fields.maxLimit, tt.fields.films...)
			got, err := r.ListDirectors(tt.args.offset, tt.args.limit)
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestFilmRepo_ListFilms(t *testing.T) {
	type fields struct {
		films    []model.Film
		maxLimit int
	}
	type args struct {
		offset int
		limit  int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []model.Film
		wantErr bool
	}{
		{
			name: "empty",
			fields: fields{
				maxLimit: 10,
			},
			args: args{
				offset: 0,
				limit:  10,
			},
		},
		{
			name: "single film",
			fields: fields{
				films: []model.Film{
					{Title: "foo", Director: "bar"},
				},
				maxLimit: 10,
			},
			args: args{
				offset: 0,
				limit:  10,
			},
			want: []model.Film{
				{Title: "foo", Director: "bar"},
			},
		},
		{
			name: "films are in alphabetical order",
			fields: fields{
				films: []model.Film{
					{Title: "foo", Director: "bar"},
					{Title: "baz", Director: "quix"},
				},
				maxLimit: 10,
			},
			args: args{
				offset: 0,
				limit:  10,
			},
			want: []model.Film{
				{Title: "baz", Director: "quix"},
				{Title: "foo", Director: "bar"},
			},
		},
		{
			name: "directors with multiple films",
			fields: fields{
				films: []model.Film{
					{Title: "foo", Director: "bar"},
					{Title: "baz", Director: "bar"},
				},
				maxLimit: 10,
			},
			args: args{
				offset: 0,
				limit:  10,
			},
			want: []model.Film{
				{Title: "baz", Director: "bar"},
				{Title: "foo", Director: "bar"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewFilmRepo(tt.fields.maxLimit, tt.fields.films...)
			got, err := r.ListFilms(tt.args.offset, tt.args.limit)
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestFilmRepo_Truncate(t *testing.T) {
	type fields struct {
		films    []model.Film
		maxLimit int
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "empty",
			fields: fields{
				maxLimit: 10,
			},
		},
		{
			name: "empty",
			fields: fields{
				films: []model.Film{
					{Title: "foo", Director: "bar"},
					{Title: "baz", Director: "bar"},
				},
				maxLimit: 10,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewFilmRepo(tt.fields.maxLimit, tt.fields.films...)

			r.Truncate()
			assert.Equal(t, 0, r.CountFilms())
			assert.Equal(t, 0, r.CountFilms())
		})
	}
}
