CREATE TABLE "movie" (
    "id" bigserial PRIMARY KEY,
    "title" varchar NOT NULL,
    "genre" varchar NOT NULL,
    "release_date" date NOT NULL,
    "created_at" timestamptz NOT NULL DEFAULT (now()),
    "updated_at" timestamptz NOT NULL DEFAULT (now())
);
