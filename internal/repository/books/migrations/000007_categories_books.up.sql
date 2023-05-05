CREATE TABLE IF NOT EXISTS categories_books(
	id SERIAL PRIMARY KEY,
	category_id INTEGER REFERENCES categories (id),
	book_id INTEGER REFERENCES books (id)
);
