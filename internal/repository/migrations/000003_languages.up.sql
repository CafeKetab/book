CREATE TABLE IF NOT EXISTS languages(
	id SERIAL PRIMARY KEY,
	name VARCHAR(30)
);

CREATE INDEX languages_name_idx ON languages (name);
