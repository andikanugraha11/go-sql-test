package psql

import (
	"testing"
	"time"

	"github.com/andikanugraha11/go-sql-test/internal/entity"
	"github.com/andikanugraha11/go-sql-test/pkg/utils"
	"github.com/stretchr/testify/require"
)

func createMovie(t *testing.T) int {
	arg := entity.Movie{
		Title:       utils.RandomString(10),
		Genre:       utils.RandomString(10),
		ReleaseDate: time.Now(),
	}

	id, err := repo.Create(&arg)

	require.NoError(t, err)
	require.NotEmpty(t, id)

	return id
}

func selectMovie(t *testing.T, id int) *entity.Movie {
	movie, err := repo.FindByID(id)
	require.NoError(t, err)
	return movie
}

func TestCreate(t *testing.T) {
	createMovie(t)
}

func TestFindByID(t *testing.T) {
	id := createMovie(t)
	movie := selectMovie(t, id)
	require.Equal(t, id, movie.ID)
}

func TestUpdate(t *testing.T) {
	id := createMovie(t)

	arg := entity.Movie{
		ID:    id,
		Title: utils.RandomString(10),
		Genre: utils.RandomString(10),
	}

	err := repo.Update(&arg)
	require.NoError(t, err)

	movie := selectMovie(t, id)

	require.NoError(t, err)
	require.Equal(t, arg.Title, movie.Title)
	require.Equal(t, arg.Genre, movie.Genre)
}

func TestDelete(t *testing.T) {
	id := createMovie(t)

	err := repo.Delete(id)
	require.NoError(t, err)

	movie := selectMovie(t, id)
	require.Nil(t, movie)
}
