postgres:
	docker-compose up -d postgres

createdb:
	docker exec -it postgres12-local-test createdb --username=root --owner=root movie_db

dropdb:
	docker exec -it postgres12-local-test dropdb movie_db

migrate-up:
	migrate -path internal/repository/psql/migration -database "postgresql://root:secret@localhost:5432/movie_db?sslmode=disable" -verbose up

migrate-down:
	migrate -path internal/repository/psql/migration -database "postgresql://root:secret@localhost:5432/movie_db?sslmode=disable" -verbose down

.PHONY: postgres createdb dropdb migrateup migratedown