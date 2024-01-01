package repository_test

import (
	"io"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/peteraba/go-htmx-playground/pkg/movies/model"
	"github.com/peteraba/go-htmx-playground/pkg/movies/repository"
)

func TestMovieRepo_CountDirectors(t *testing.T) {
	t.Parallel()

	dummyLogger := slog.New(slog.NewTextHandler(io.Discard, nil))

	type fields struct {
		movies   []model.Movie
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
				movies:   []model.Movie{},
				maxLimit: 10,
			},
			want: 0,
		},
		{
			name: "single movie",
			fields: fields{
				movies: []model.Movie{
					{Title: "Forrest Gump", Director: "Ethan White"},
				},
				maxLimit: 10,
			},
			want: 1,
		},
		{
			name: "no directors with multiple movies",
			fields: fields{
				movies: []model.Movie{
					{Title: "Forrest Gump", Director: "Ethan White"},
					{Title: "Die Hard", Director: "Kara Nader"},
				},
				maxLimit: 10,
			},
			want: 2,
		},
		{
			name: "directors with multiple movies",
			fields: fields{
				movies: []model.Movie{
					{Title: "Forrest Gump", Director: "Ethan White"},
					{Title: "Die Hard", Director: "Ethan White"},
				},
				maxLimit: 10,
			},
			want: 1,
		},
	}
	for _, ttt := range tests {
		tt := ttt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			r := repository.NewMovieRepo(dummyLogger, tt.fields.maxLimit, tt.fields.movies...)
			if got := r.CountDirectors(); got != tt.want {
				t.Errorf("CountDirectors() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMovieRepo_CountMovies(t *testing.T) {
	t.Parallel()

	dummyLogger := slog.New(slog.NewTextHandler(io.Discard, nil))

	type fields struct {
		movies   []model.Movie
		maxLimit int
	}
	type args struct {
		searchTerm string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   int
	}{
		{
			name: "empty",
			fields: fields{
				maxLimit: 10,
				movies:   nil,
			},
			args: args{
				searchTerm: "",
			},
			want: 0,
		},
		{
			name: "single movie",
			fields: fields{
				movies: []model.Movie{
					{Title: "Forrest Gump", Director: "Ethan White"},
				},
				maxLimit: 10,
			},
			args: args{
				searchTerm: "",
			},
			want: 1,
		},
		{
			name: "no directors with multiple movies",
			fields: fields{
				movies: []model.Movie{
					{Title: "Forrest Gump", Director: "Ethan White"},
					{Title: "Die Hard", Director: "Kara Nader"},
				},
				maxLimit: 10,
			},
			args: args{
				searchTerm: "",
			},
			want: 2,
		},
		{
			name: "directors with multiple movies",
			fields: fields{
				movies: []model.Movie{
					{Title: "Forrest Gump", Director: "Ethan White"},
					{Title: "Die Hard", Director: "Ethan White"},
				},
				maxLimit: 10,
			},
			args: args{
				searchTerm: "",
			},
			want: 2,
		},
		{
			name: "directors with multiple movies and search term",
			fields: fields{
				movies: []model.Movie{
					{Title: "Forrest Gump", Director: "Ethan White"},
					{Title: "Die Hard", Director: "Ethan White"},
				},
				maxLimit: 10,
			},
			args: args{
				searchTerm: "orrest",
			},
			want: 1,
		},
	}
	for _, ttt := range tests {
		tt := ttt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			r := repository.NewMovieRepo(dummyLogger, tt.fields.maxLimit, tt.fields.movies...)
			if got := r.CountMovies(tt.args.searchTerm); got != tt.want {
				t.Errorf("CountMovies() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMovieRepo_Insert(t *testing.T) {
	t.Parallel()

	dummyLogger := slog.New(slog.NewTextHandler(io.Discard, nil))

	type fields struct {
		movies   []model.Movie
		maxLimit int
	}
	type args struct {
		newMovies []model.Movie
	}
	tests := []struct {
		name              string
		fields            fields
		args              args
		wantMovieCount    int
		wantDirectorCount int
	}{
		{
			name: "from empty",
			fields: fields{
				maxLimit: 10,
				movies:   nil,
			},
			args: args{
				newMovies: []model.Movie{
					{Title: "Forrest Gump", Director: "Ethan White"},
					{Title: "Die Hard", Director: "Ethan White"},
				},
			},
			wantMovieCount:    2,
			wantDirectorCount: 1,
		},
		{
			name: "from non-empty",
			fields: fields{
				maxLimit: 10,
				movies: []model.Movie{
					{Title: "Back to the Future", Director: "Ethan White"},
				},
			},
			args: args{
				newMovies: []model.Movie{
					{Title: "Forrest Gump", Director: "Ethan White"},
					{Title: "Die Hard", Director: "Ethan White"},
				},
			},
			wantMovieCount:    3,
			wantDirectorCount: 1,
		},
		{
			name: "duplicate movieKeys get skipped",
			fields: fields{
				maxLimit: 10,
				movies: []model.Movie{
					{Title: "Forrest Gump", Director: "Kara Nader"},
				},
			},
			args: args{
				newMovies: []model.Movie{
					{Title: "Forrest Gump", Director: "Ethan White"},
					{Title: "Die Hard", Director: "Ethan White"},
				},
			},
			wantMovieCount:    2,
			wantDirectorCount: 2,
		},
	}
	for _, ttt := range tests {
		tt := ttt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			r := repository.NewMovieRepo(dummyLogger, tt.fields.maxLimit, tt.fields.movies...)
			_ = r.Insert(tt.args.newMovies...)

			assert.Equal(t, tt.wantMovieCount, r.CountMovies(""))
			assert.Equal(t, tt.wantDirectorCount, r.CountDirectors())
		})
	}
}

func TestMovieRepo_ListDirectors(t *testing.T) {
	t.Parallel()

	dummyLogger := slog.New(slog.NewTextHandler(io.Discard, nil))

	type fields struct {
		movies   []model.Movie
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
				movies:   nil,
			},
			args: args{
				offset: 0,
				limit:  10,
			},
			want: nil,
		},
		{
			name: "single movie",
			fields: fields{
				movies: []model.Movie{
					{Title: "Forrest Gump", Director: "Ethan White"},
				},
				maxLimit: 10,
			},
			args: args{
				offset: 0,
				limit:  10,
			},
			want: []model.Director{
				{Name: "Ethan White", Titles: []string{"Forrest Gump"}},
			},
		},
		{
			name: "no directors with multiple movies",
			fields: fields{
				movies: []model.Movie{
					{Title: "Forrest Gump", Director: "Ethan White"},
					{Title: "Die Hard", Director: "Kara Nader"},
				},
				maxLimit: 10,
			},
			args: args{
				offset: 0,
				limit:  10,
			},
			want: []model.Director{
				{Name: "Ethan White", Titles: []string{"Forrest Gump"}},
				{Name: "Kara Nader", Titles: []string{"Die Hard"}},
			},
		},
		{
			name: "directors with multiple movies",
			fields: fields{
				movies: []model.Movie{
					{Title: "Forrest Gump", Director: "Ethan White"},
					{Title: "Die Hard", Director: "Ethan White"},
				},
				maxLimit: 10,
			},
			args: args{
				offset: 0,
				limit:  10,
			},
			want: []model.Director{
				{Name: "Ethan White", Titles: []string{"Die Hard", "Forrest Gump"}},
			},
		},
	}
	for _, ttt := range tests {
		tt := ttt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			r := repository.NewMovieRepo(dummyLogger, tt.fields.maxLimit, tt.fields.movies...)
			got, err := r.ListDirectors(tt.args.offset, tt.args.limit)
			require.NoError(t, err)

			assert.Equal(t, tt.want, got)
		})
	}
}

func TestMovieRepo_ListMovies(t *testing.T) {
	t.Parallel()

	dummyLogger := slog.New(slog.NewTextHandler(io.Discard, nil))

	type fields struct {
		movies   []model.Movie
		maxLimit int
	}

	type args struct {
		offset     int
		limit      int
		searchTerm string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []model.Movie
	}{
		{
			name: "empty",
			fields: fields{
				maxLimit: 10,
				movies:   nil,
			},
			args: args{
				offset:     0,
				limit:      10,
				searchTerm: "",
			},
			want: []model.Movie{},
		},
		{
			name: "single movie",
			fields: fields{
				movies: []model.Movie{
					{Title: "Forrest Gump", Director: "Ethan White"},
				},
				maxLimit: 10,
			},
			args: args{
				offset:     0,
				limit:      10,
				searchTerm: "",
			},
			want: []model.Movie{
				{Title: "Forrest Gump", Director: "Ethan White"},
			},
		},
		{
			name: "movies are in alphabetical order",
			fields: fields{
				movies: []model.Movie{
					{Title: "Forrest Gump", Director: "Ethan White"},
					{Title: "Die Hard", Director: "Kara Nader"},
				},
				maxLimit: 10,
			},
			args: args{
				offset:     0,
				limit:      10,
				searchTerm: "",
			},
			want: []model.Movie{
				{Title: "Die Hard", Director: "Kara Nader"},
				{Title: "Forrest Gump", Director: "Ethan White"},
			},
		},
		{
			name: "directors with multiple movies",
			fields: fields{
				movies: []model.Movie{
					{Title: "Forrest Gump", Director: "Ethan White"},
					{Title: "Die Hard", Director: "Ethan White"},
				},
				maxLimit: 10,
			},
			args: args{
				offset:     0,
				limit:      10,
				searchTerm: "",
			},
			want: []model.Movie{
				{Title: "Die Hard", Director: "Ethan White"},
				{Title: "Forrest Gump", Director: "Ethan White"},
			},
		},
		{
			name: "filtering works",
			fields: fields{
				movies: []model.Movie{
					{Title: "Forrest Gump", Director: "Ethan White"},
					{Title: "Die Hard", Director: "Ethan White"},
				},
				maxLimit: 10,
			},
			args: args{
				offset:     0,
				limit:      10,
				searchTerm: "orre",
			},
			want: []model.Movie{
				{Title: "Forrest Gump", Director: "Ethan White"},
			},
		},
		{
			name: "directors with multiple movies can be filtered out",
			fields: fields{
				movies: []model.Movie{
					{Title: "Forrest Gump", Director: "Ethan White"},
					{Title: "Die Hard", Director: "Ethan White"},
				},
				maxLimit: 10,
			},
			args: args{
				offset:     0,
				limit:      10,
				searchTerm: "nope",
			},
			want: []model.Movie{},
		},
	}
	for _, ttt := range tests {
		tt := ttt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			r := repository.NewMovieRepo(dummyLogger, tt.fields.maxLimit, tt.fields.movies...)
			got, err := r.ListMovies(tt.args.offset, tt.args.limit, tt.args.searchTerm)
			require.NoError(t, err)

			assert.Equal(t, tt.want, got)
		})
	}
}

func TestMovieRepo_Truncate(t *testing.T) {
	t.Parallel()

	dummyLogger := slog.New(slog.NewTextHandler(io.Discard, nil))

	type fields struct {
		movies   []model.Movie
		maxLimit int
	}
	type args struct {
		searchTerm string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "empty",
			fields: fields{
				maxLimit: 10,
				movies:   nil,
			},
			args: args{
				searchTerm: "",
			},
		},
		{
			name: "empty",
			fields: fields{
				movies: []model.Movie{
					{Title: "Forrest Gump", Director: "Ethan White"},
					{Title: "Die Hard", Director: "Ethan White"},
				},
				maxLimit: 10,
			},
			args: args{
				searchTerm: "",
			},
		},
	}
	for _, ttt := range tests {
		tt := ttt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			r := repository.NewMovieRepo(dummyLogger, tt.fields.maxLimit, tt.fields.movies...)

			r.Truncate()
			assert.Equal(t, 0, r.CountMovies(tt.args.searchTerm))
		})
	}
}

func TestMovieRepo_DeleteByKey(t *testing.T) {
	t.Parallel()

	const maxLimit = 10

	dummyLogger := slog.New(slog.NewTextHandler(io.Discard, nil))

	type fields struct {
		movies []model.Movie
	}
	type args struct {
		keys []string
	}
	type want struct {
		movieTitles   []string
		directorNames []string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "delete from empty",
			fields: fields{
				movies: []model.Movie{},
			},
			args: args{
				keys: []string{"forrest-gump"},
			},
			want: want{
				movieTitles:   []string{},
				directorNames: []string{},
			},
		},
		{
			name: "only keys to be deleted",
			fields: fields{
				movies: []model.Movie{
					{Title: "Forrest Gump", Director: "Ethan White"},
				},
			},
			args: args{
				keys: []string{"forrest-gump"},
			},
			want: want{
				movieTitles:   []string{},
				directorNames: []string{},
			},
		},
		{
			name: "first keys to be deleted",
			fields: fields{
				movies: []model.Movie{
					{Title: "Forrest Gump", Director: "Ethan White"},
					{Title: "Die Hard", Director: "Ethan White"},
				},
			},
			args: args{
				keys: []string{"forrest-gump"},
			},
			want: want{
				movieTitles:   []string{"Die Hard"},
				directorNames: []string{"Ethan White"},
			},
		},
		{
			name: "second keys to be deleted",
			fields: fields{
				movies: []model.Movie{
					{Title: "Forrest Gump", Director: "Ethan White"},
					{Title: "Die Hard", Director: "Ethan White"},
				},
			},
			args: args{
				keys: []string{"die-hard"},
			},
			want: want{
				movieTitles:   []string{"Forrest Gump"},
				directorNames: []string{"Ethan White"},
			},
		},
		{
			name: "middle keys to be deleted",
			fields: fields{
				movies: []model.Movie{
					{Title: "Forrest Gump", Director: "Ethan White"},
					{Title: "Die Hard", Director: "Ethan White"},
					{Title: "Fight Club", Director: "Shaylee Hegmann"},
				},
			},
			args: args{
				keys: []string{"die-hard"},
			},
			want: want{
				movieTitles:   []string{"Fight Club", "Forrest Gump"},
				directorNames: []string{"Ethan White", "Shaylee Hegmann"},
			},
		},
	}
	for _, ttt := range tests {
		tt := ttt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			sut := repository.NewMovieRepo(dummyLogger, maxLimit, tt.fields.movies...)

			sut.DeleteMoviesByKey(tt.args.keys...)

			gotTitles := sut.ListAllTitles()
			assert.Equal(t, tt.want.movieTitles, gotTitles)

			gotDirectorNames := sut.ListAllDirectorNames()
			assert.Equal(t, tt.want.directorNames, gotDirectorNames)
		})
	}
}
