package psql

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
)

const (
	dbDriver = "postgres"
	dbSource = "postgresql://test:test@localhost:%s/movie_db?sslmode=disable"
)

var (
	testDB *sql.DB
	repo   *MovieRepository
)

func TestMain(m *testing.M) {
	// Set the timeout duration for test setup tasks
	setupTimeoutDuration := 5 * time.Minute

	// Use a channel to receive the setup task completion signal
	setupDone := make(chan bool)

	var (
		pool     *dockertest.Pool
		resource *dockertest.Resource
	)
	// Start a goroutine to perform the setup tasks
	go func() {
		var err error
		pool, err = dockertest.NewPool("")
		if err != nil {
			log.Fatalf("Could not construct pool: %s", err)
		}
		err = pool.Client.Ping()
		if err != nil {
			log.Fatalf("Could not connect to Docker: %s", err)
		}

		resource, err = pool.RunWithOptions(&dockertest.RunOptions{
			Repository: "postgres",
			Tag:        "12",
			Env: []string{
				"POSTGRES_USER=test",
				"POSTGRES_PASSWORD=test",
				"POSTGRES_DB=movie_db",
				"listen_addresses = '*'",
			},
		}, func(config *docker.HostConfig) {
			// set AutoRemove to true so that stopped container goes away by itself
			config.AutoRemove = true
			config.RestartPolicy = docker.RestartPolicy{
				Name: "no",
			}
		})
		// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
		if err := pool.Retry(func() error {
			var err error
			testDB, err = sql.Open(dbDriver, fmt.Sprintf(dbSource, resource.GetPort("5432/tcp")))
			if err != nil {
				return err
			}

			return testDB.Ping()
		}); err != nil {
			log.Fatalf("Could not connect to database: %s", err)
		}

		driver, err := postgres.WithInstance(testDB, &postgres.Config{})
		if err != nil {
			log.Fatalf("Could not migrate database: %s", err)
		}

		migration, err := migrate.NewWithDatabaseInstance(
			"file://"+filepath.Join(".", "migrations"),
			"postgres", driver,
		)
		if err != nil {
			log.Fatalf("Failed to initialize migration instance: %s", err)
		}

		err = migration.Up()
		if err != nil {
			log.Fatalf("Failed to apply migrations: %s", err)
		}

		setupDone <- true
	}()

	// Wait for either the setup tasks to complete or the timeout duration to be reached
	select {
	case <-setupDone:
		// Setup tasks completed within the timeout duration
	case <-time.After(setupTimeoutDuration):
		// Setup tasks timed out
		fmt.Println("Test setup tasks timed out")
		return
	}

	// prepareTestingEnvirontment(m)
	repo = NewMovieRepository(testDB)
	code := m.Run()
	// You can't defer this because os.Exit doesn't care for defer
	if err := pool.Purge(resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}

	os.Exit(code)
}
