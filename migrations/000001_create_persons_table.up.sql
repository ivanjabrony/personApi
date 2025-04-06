CREATE TABLE persons (
  id SERIAL PRIMARY KEY,
  name TEXT NOT NULL,
  surname TEXT NOT NULL,
  patronymic TEXT NULL,
  age INT NULL,
  gender TEXT NULL,
  nationality TEXT NULL
);

CREATE INDEX on persons(id);