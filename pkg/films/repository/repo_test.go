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
					{Title: "Forrest Gump", Director: "Ethan White"},
				},
				maxLimit: 10,
			},
			want: 1,
		},
		{
			name: "no directors with multiple films",
			fields: fields{
				films: []model.Film{
					{Title: "Forrest Gump", Director: "Ethan White"},
					{Title: "Die Hard", Director: "Kara Nader"},
				},
				maxLimit: 10,
			},
			want: 2,
		},
		{
			name: "directors with multiple films",
			fields: fields{
				films: []model.Film{
					{Title: "Forrest Gump", Director: "Ethan White"},
					{Title: "Die Hard", Director: "Ethan White"},
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
					{Title: "Forrest Gump", Director: "Ethan White"},
				},
				maxLimit: 10,
			},
			want: 1,
		},
		{
			name: "no directors with multiple films",
			fields: fields{
				films: []model.Film{
					{Title: "Forrest Gump", Director: "Ethan White"},
					{Title: "Die Hard", Director: "Kara Nader"},
				},
				maxLimit: 10,
			},
			want: 2,
		},
		{
			name: "directors with multiple films",
			fields: fields{
				films: []model.Film{
					{Title: "Forrest Gump", Director: "Ethan White"},
					{Title: "Die Hard", Director: "Ethan White"},
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
					{Title: "Forrest Gump", Director: "Ethan White"},
					{Title: "Die Hard", Director: "Ethan White"},
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
					{Title: "Back to the Future", Director: "Ethan White"},
				},
			},
			args: args{
				newFilms: []model.Film{
					{Title: "Forrest Gump", Director: "Ethan White"},
					{Title: "Die Hard", Director: "Ethan White"},
				},
			},
			wantFilmCount:     3,
			wantDirectorCount: 1,
		},
		{
			name: "duplicate filmTitles get skipped",
			fields: fields{
				maxLimit: 10,
				films: []model.Film{
					{Title: "Forrest Gump", Director: "Kara Nader"},
				},
			},
			args: args{
				newFilms: []model.Film{
					{Title: "Forrest Gump", Director: "Ethan White"},
					{Title: "Die Hard", Director: "Ethan White"},
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
			name: "no directors with multiple films",
			fields: fields{
				films: []model.Film{
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
			name: "directors with multiple films",
			fields: fields{
				films: []model.Film{
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
					{Title: "Forrest Gump", Director: "Ethan White"},
				},
				maxLimit: 10,
			},
			args: args{
				offset: 0,
				limit:  10,
			},
			want: []model.Film{
				{Title: "Forrest Gump", Director: "Ethan White"},
			},
		},
		{
			name: "films are in alphabetical order",
			fields: fields{
				films: []model.Film{
					{Title: "Forrest Gump", Director: "Ethan White"},
					{Title: "Die Hard", Director: "Kara Nader"},
				},
				maxLimit: 10,
			},
			args: args{
				offset: 0,
				limit:  10,
			},
			want: []model.Film{
				{Title: "Die Hard", Director: "Kara Nader"},
				{Title: "Forrest Gump", Director: "Ethan White"},
			},
		},
		{
			name: "directors with multiple films",
			fields: fields{
				films: []model.Film{
					{Title: "Forrest Gump", Director: "Ethan White"},
					{Title: "Die Hard", Director: "Ethan White"},
				},
				maxLimit: 10,
			},
			args: args{
				offset: 0,
				limit:  10,
			},
			want: []model.Film{
				{Title: "Die Hard", Director: "Ethan White"},
				{Title: "Forrest Gump", Director: "Ethan White"},
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
					{Title: "Forrest Gump", Director: "Ethan White"},
					{Title: "Die Hard", Director: "Ethan White"},
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

func TestFilmRepo_DeleteByTitle(t *testing.T) {
	const maxLimit = 10

	type fields struct {
		films []model.Film
	}
	type args struct {
		title string
	}
	type want struct {
		filmTitles    []string
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
				films: []model.Film{},
			},
			args: args{
				title: "Forrest Gump",
			},
			want: want{
				filmTitles:    []string{},
				directorNames: []string{},
			},
		},
		{
			name: "only title to be deleted",
			fields: fields{
				films: []model.Film{
					{Title: "Forrest Gump", Director: "Ethan White"},
				},
			},
			args: args{
				title: "Forrest Gump",
			},
			want: want{
				filmTitles:    []string{},
				directorNames: []string{},
			},
		},
		{
			name: "first title to be deleted",
			fields: fields{
				films: []model.Film{
					{Title: "Forrest Gump", Director: "Ethan White"},
					{Title: "Die Hard", Director: "Ethan White"},
				},
			},
			args: args{
				title: "Forrest Gump",
			},
			want: want{
				filmTitles:    []string{"Die Hard"},
				directorNames: []string{"Ethan White"},
			},
		},
		{
			name: "second title to be deleted",
			fields: fields{
				films: []model.Film{
					{Title: "Forrest Gump", Director: "Ethan White"},
					{Title: "Die Hard", Director: "Ethan White"},
				},
			},
			args: args{
				title: "Die Hard",
			},
			want: want{
				filmTitles:    []string{"Forrest Gump"},
				directorNames: []string{"Ethan White"},
			},
		},
		{
			name: "middle title to be deleted",
			fields: fields{
				films: []model.Film{
					{Title: "Forrest Gump", Director: "Ethan White"},
					{Title: "Die Hard", Director: "Ethan White"},
					{Title: "Fight Club", Director: "Shaylee Hegmann"},
				},
			},
			args: args{
				title: "Die Hard",
			},
			want: want{
				filmTitles:    []string{"Fight Club", "Forrest Gump"},
				directorNames: []string{"Ethan White", "Shaylee Hegmann"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewFilmRepo(maxLimit, tt.fields.films...)

			r.DeleteByTitle(tt.args.title)

			assert.Equal(t, tt.want.filmTitles, r.filmTitles)
			assert.Equal(t, tt.want.directorNames, r.directorNames)
		})
	}
}
