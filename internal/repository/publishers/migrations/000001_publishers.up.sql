CREATE TABLE IF NOT EXISTS publishers(
	id SERIAL PRIMARY KEY,
	name VARCHAR(30),
	title VARCHAR(80),
	description VARCHAR(255)
);

CREATE INDEX publishers_name_idx ON publishers (name);
