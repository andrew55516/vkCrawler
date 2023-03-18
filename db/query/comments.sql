-- name: WriteDownComment :one
INSERT INTO comments (post_id,
                      owner)
VALUES ($1, $2) RETURNING *;

-- name: GetComment :one
SELECT * from comments
WHERE post_id = $1 AND owner = $2;