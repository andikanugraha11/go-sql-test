package psql

import (
	"database/sql"

	"github.com/andikanugraha11/go-sql-test/internal/entity"
)

type MovieRepository struct {
	db *sql.DB
}

func NewMovieRepository(db *sql.DB) *MovieRepository {
	return &MovieRepository{
		db: db,
	}
}

func (r *MovieRepository) FindByID(id int) (*entity.Movie, error) {
	var movie entity.Movie

	query := `
        SELECT id, title, genre, release_date, created_at, updated_at FROM movie WHERE id = $1
    `
	err := r.db.QueryRow(query, id).Scan(&movie.ID, &movie.Title, &movie.Genre, &movie.ReleaseDate, &movie.CreatedAt, &movie.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &movie, nil
}

func (r *MovieRepository) Create(movie *entity.Movie) (int, error) {
	var id int
	query := `
        INSERT INTO movie (title, genre, release_date) VALUES ($1, $2, $3) RETURNING id
    `
	err := r.db.QueryRow(query, movie.Title, movie.Genre, movie.ReleaseDate).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *MovieRepository) Update(movie *entity.Movie) error {
	query := `
        UPDATE movie SET title = $1, genre = $2 WHERE id = $3
    `
	_, err := r.db.Exec(query, movie.Title, movie.Genre, movie.ID)
	if err != nil {
		return err
	}

	return nil
}

func (r *MovieRepository) Delete(id int) error {
	query := `
        DELETE FROM movie WHERE id = $1
    `
	_, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	return nil
}
