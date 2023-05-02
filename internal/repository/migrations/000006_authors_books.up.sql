CREATE TABLE IF NOT EXISTS authors_books(
	id SERIAL PRIMARY KEY,
	author_id INTEGER REFERENCES authors (id),
	book_id INTEGER REFERENCES books (id)
);
