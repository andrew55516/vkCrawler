-- name: WriteDownComment :one
INSERT INTO comments (post_id,
                      owner,
                      thread_owner,
                      created_at)
VALUES ($1, $2, $3, $4) RETURNING *;

-- name: GetComment :one
SELECT * from comments
WHERE post_id = $1 AND owner = $2;