-- +goose Up
-- +goose StatementBegin
CREATE TABLE "actors"
(
    "id"        uuid PRIMARY KEY,
    "name"      varchar,
    "surname"   varchar,
    "birthdate" timestamp,
    "movie_id"   uuid,
    "sex"       varchar
);

CREATE TABLE "users"
(
    "id"         uuid PRIMARY KEY,
    "login"   varchar,
    "password" varchar,
    "role"       varchar,
    "created_at" timestamp
);

CREATE TABLE "movies"
(
    "id"          UUID PRIMARY KEY,
    "title"       varchar NOT NULL,
    "description" varchar(1000),
    "rating"      real CHECK (rating >= 0 AND rating <= 10),
    "created_at"  timestamp,
    CONSTRAINT length_constraint_title CHECK (LENGTH("title") >= 1 AND LENGTH("title") <= 150)
);

CREATE TABLE "actors_movies"
(
    "actors_movie_id" uuid,
    "movies_id"       UUID,
    PRIMARY KEY ("actors_movie_id", "movies_id"),
    FOREIGN KEY ("actors_movie_id") REFERENCES "actors" ("id") ON DELETE CASCADE,
    FOREIGN KEY ("movies_id") REFERENCES movies ("id") ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS "actors_movies";
DROP TABLE IF EXISTS "users";
DROP TABLE IF EXISTS "actors";
DROP TABLE IF EXISTS movies;
-- +goose StatementEnd