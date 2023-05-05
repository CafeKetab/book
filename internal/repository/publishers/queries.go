package publishers

const QueryInsert = `
INSERT INTO 
	publishers(name, title, description) 
	VALUES($1, $2, $3) 
RETURNING id;`

const QueryGetDetail = `
SELECT name, title, description 
FROM publishers 
WHERE id=$1;`

const QueryGetAll = `
SELECT id, name, title 
FROM publishers 
WHERE 
	id > $1 AND
	name LIKE '%' || $2 || '%'
ORDER BY id
FETCH NEXT $3 ROWS ONLY;`
