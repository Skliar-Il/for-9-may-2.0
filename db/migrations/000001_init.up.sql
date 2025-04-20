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
    rank text NOT NULL
);

CREATE TABLE IF NOT EXISTS form (
    id uuid DEFAULT uuid_generate_v4() PRIMARY KEY,
    status_check boolean NOT NULL DEFAULT false,
    date_published timestamp NOT NULL DEFAULT TIMEZONE('utc', now()),
    person_id uuid NOT NULL UNIQUE REFERENCES person(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS owner (
    id uuid DEFAULT uuid_generate_v4() PRIMARY KEY,
    email text NOT NULL,
    telegram text NOT NULL,
    name text NOT NULL,
    surname text NOT NULL,
    patronymic text,
    relative text NOT NULL,
    form_id uuid NOT NULL UNIQUE REFERENCES form(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS person_photo (
    id serial PRIMARY KEY,
    person_id uuid NOT NULL REFERENCES person(id) ON DELETE CASCADE ,
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
    person_id uuid NOT NULL REFERENCES person(id) ON DELETE CASCADE ,
    medal_id int NOT NULL REFERENCES medal (id) ON DELETE CASCADE ,
    UNIQUE (person_id, medal_id)
);

CREATE OR REPLACE FUNCTION check_main_photo_before_insert()
    RETURNS TRIGGER AS $$
BEGIN
    -- Если пытаемся добавить НЕ главное фото, но у пользователя ещё нет главного фото
    IF NEW.main_status = FALSE AND NOT EXISTS (
        SELECT 1 FROM person_photo
        WHERE person_id = NEW.person_id AND main_status = TRUE
    ) THEN
        RAISE EXCEPTION 'Cannot add non-main photo: person must have at least one main photo first';
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER prevent_non_main_photo_without_main
    BEFORE INSERT ON person_photo
    FOR EACH ROW
EXECUTE FUNCTION check_main_photo_before_insert();

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
    o.name AS owner_name,
    o.surname AS owner_surname,
    o.patronymic AS owner_patronymic,
    o.telegram,
    o.relative,
    f.status_check AS status_check,
    f.date_published,
    (
        SELECT COALESCE(JSON_AGG(JSON_BUILD_OBJECT(
                'id', m.id,
                'name', m.name,
                'photo_link', m.photo_link
                                 )), '[]')
        FROM medal_person mp
        JOIN medal m ON m.id = mp.medal_id
        WHERE mp.person_id = p.id
    ) AS medals,
    (
        SELECT COALESCE(JSON_AGG(JSON_BUILD_OBJECT(
                                'id', pp.id,
                                 'link', pp.link,
                                 'is_main', pp.main_status
                                 )), '[]')
        FROM person_photo pp
        WHERE pp.person_id = p.id
    ) AS photo
    FROM person p
    LEFT JOIN form f ON f.person_id = p.id
    LEFT JOIN owner o ON o.form_id = f.id;

