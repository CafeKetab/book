package categories

const QueryInsert = `
INSERT INTO 
	languages(name) 
	VALUES($1) 
RETURNING id;`

const QueryGetDetail = `
SELECT name 
FROM languages 
WHERE id=$1;`

const QueryGetAll = `
SELECT id, name 
FROM languages 
WHERE 
	id > $1 AND
	name LIKE '%' || $2 || '%'
ORDER BY id
FETCH NEXT $3 ROWS ONLY;`
