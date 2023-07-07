package psql

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

const (
	dbDriver = "postgres"
	dbSource = "postgresql://root:secret@localhost:5432/movie_db?sslmode=disable"
)

var (
	testQueries *sql.DB
	repo        *MovieRepository
)

func TestMain(m *testing.M) {
	var err error
	testQueries, err = sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("Cannot connect to database", err)
	}

	repo = NewMovieRepository(testQueries)

	os.Exit(m.Run())
}
