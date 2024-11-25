-- +migrate Up
CREATE TABLE songs (
    id serial primary key,
    "group" text not null,
    song text not null,
    release_date timestamp not null,
    verses text [] not null,
    link text not null
);
-- +migrate Down
DROP TABLE songs;