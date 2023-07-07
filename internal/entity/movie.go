package entity

import "time"

type Movie struct {
	ID          int       `db:"id"`
	Title       string    `db:"title"`
	Genre       string    `db:"genre"`
	ReleaseDate time.Time `db:"release_date"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}
