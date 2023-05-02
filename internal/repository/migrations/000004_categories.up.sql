CREATE TABLE IF NOT EXISTS categories(
	id SERIAL PRIMARY KEY,
	name VARCHAR(30),
	description VARCHAR(255)
);

CREATE INDEX categories_name_idx ON categories (name);
