package psql

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

const (
	dbDriver   = "postgres"
	dbSource   = "postgresql://test:test@localhost:%s/movie_db?sslmode=disable"
	migrations = "file://migrations"
)

var (
	testDB   *sql.DB
	repo     *MovieRepository
	pool     *dockertest.Pool
	resource *dockertest.Resource
)

func TestMain(m *testing.M) {
	// Set up the Docker test environment
	setupDockerTestEnvironment()

	// Apply database migrations
	applyDatabaseMigrations()

	repo = NewMovieRepository(testDB)

	// Run the tests
	code := m.Run()

	// Tear down the Docker test environment
	tearDownDockerTestEnvironment()

	// Exit with the appropriate exit code
	os.Exit(code)
}

func setupDockerTestEnvironment() {
	setupTimeoutDuration := 5 * time.Minute
	setupDone := make(chan bool)

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
			config.AutoRemove = true
			config.RestartPolicy = docker.RestartPolicy{
				Name: "no",
			}
		})

		if err != nil {
			log.Fatalf("Failed to create resource: %s", err)
		}

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

		setupDone <- true
	}()

	select {
	case <-setupDone:
		log.Println("Docker test environment setup completed")
	case <-time.After(setupTimeoutDuration):
		log.Println("Docker test environment setup timed out")
	}
}

func applyDatabaseMigrations() {
	driver, err := postgres.WithInstance(testDB, &postgres.Config{})
	if err != nil {
		log.Fatalf("Could not create migration driver: %s", err)
	}

	migration, err := migrate.NewWithDatabaseInstance(
		migrations,
		"postgres", driver,
	)
	if err != nil {
		log.Fatalf("Failed to initialize migration instance: %s", err)
	}

	if err := migration.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("Failed to apply migrations: %s", err)
	}

	log.Println("Database migrations applied successfully")
}

func tearDownDockerTestEnvironment() {
	if err := pool.Purge(resource); err != nil {
		log.Fatalf("Could not purge Docker resource: %s", err)
	}
}
