CREATE TABLE IF NOT EXISTS authors(
	id SERIAL PRIMARY KEY,
	full_name VARCHAR(30)
);

CREATE INDEX authors_name_idx ON authors (full_name);
