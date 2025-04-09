CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS person (
    id uuid DEFAULT uuid_generate_v4() PRIMARY KEY,
    name text NOT NULL,
    surname text NOT NULL,
    patronymic text,
    date_death int,
    date_birth int,
    city_birth text,
    history text NOT NULL,
    rank text NOT NULL,
    UNIQUE (name, surname, COALESCE(patronymic, ''))
);

CREATE TABLE IF NOT EXISTS form (
    id uuid DEFAULT uuid_generate_v4() PRIMARY KEY,
    status_check boolean NOT NULL DEFAULT false,
    date_published timestamp NOT NULL DEFAULT TIMEZONE('utc', now()),
    person_id uuid NOT NULL UNIQUE REFERENCES person(id)
);

CREATE TABLE IF NOT EXISTS owner (
    id uuid DEFAULT uuid_generate_v4() PRIMARY KEY,
    email text NOT NULL,
    telegram text NOT NULL,
    name text NOT NULL,
    surname text NOT NULL,
    patronymic text,
    relative text NOT NULL,
    form_id uuid NOT NULL UNIQUE REFERENCES form(id)
);

CREATE TABLE IF NOT EXISTS person_photo (
    id serial PRIMARY KEY,
    person_id uuid NOT NULL REFERENCES person(id),
    link text NOT NULL UNIQUE,
    main_status boolean NOT NULL DEFAULT FALSE
);
CREATE UNIQUE INDEX idx_only_one_main_photo
    ON person_photo (person_id)
    WHERE main_status = TRUE;

CREATE TABLE IF NOT EXISTS medal (
    id serial PRIMARY KEY,
    name text NOT NULL,
    photo_link text NOT NULL
);

CREATE TABLE IF NOT EXISTS medal_person (
    id serial PRIMARY KEY,
    person_id uuid NOT NULL REFERENCES person(id),
    medal_id int NOT NULL REFERENCES medal (id),
    UNIQUE (person_id, medal_id)
)