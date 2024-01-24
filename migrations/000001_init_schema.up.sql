CREATE TABLE IF NOT EXISTS persons (
                         id BIGSERIAL PRIMARY KEY,
                         name VARCHAR(100) NOT NULL,
                         surname VARCHAR(100) NOT NULL,
                         patronymic VARCHAR(200),
                         age INT NOT NULL,
                         gender VARCHAR(6) NOT NULL,
                         nationality VARCHAR(2) NOT NULL
);

CREATE INDEX ON "persons" ("id");