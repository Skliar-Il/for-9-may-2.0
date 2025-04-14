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
    UNIQUE (name, surname)
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
    name text UNIQUE NOT NULL,
    photo_link text UNIQUE
);

CREATE TABLE IF NOT EXISTS medal_person (
    id serial PRIMARY KEY,
    person_id uuid NOT NULL REFERENCES person(id),
    medal_id int NOT NULL REFERENCES medal (id),
    UNIQUE (person_id, medal_id)
);


CREATE OR REPLACE VIEW all_person_fields_view AS
SELECT
    p.id,
    p.name,
    p.surname,
    p.patronymic,
    p.date_birth,
    p.date_death,
    p.city_birth,
    p.history,
    p.rank,
    o.email,
    o.name as owner_name,
    o.surname as owner_surname,
    o.patronymic as owner_patronymic,
    o.telegram,
    o.relative,
    f.status_check as status_check,
    (
        SELECT COALESCE(JSON_AGG(JSON_BUILD_OBJECT(
                'id', m.id,
                'name', m.name,
                'url', m.photo_link
                                 )), '[]')
        FROM medal_person mp
                 JOIN medal m ON m.id = mp.medal_id
        WHERE mp.person_id = p.id
    ) AS medals
FROM person p
         LEFT JOIN form f ON f.person_id = p.id
         LEFT JOIN owner o ON o.form_id = f.id;
