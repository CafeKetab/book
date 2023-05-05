package authors

const QueryInsert = `
INSERT INTO 
	authors(full_name) 
	VALUES($1) 
RETURNING id;`

const QueryGetDetail = `
SELECT full_name
FROM authors 
WHERE id=$1;`

const QueryGetAll = `
SELECT id, full_name
FROM authors 
WHERE 
	id > $1 AND
	name LIKE '%' || $2 || '%'
ORDER BY id
FETCH NEXT $3 ROWS ONLY;`
