CREATE TABLE IF NOT EXISTS books(
	id SERIAL PRIMARY KEY,
	name VARCHAR(30),
	description VARCHAR(255),
	publisher_id INTEGER REFERENCES publishers (id),
	language_id INTEGER REFERENCES languages (id),
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX books_name_idx ON books (name);
