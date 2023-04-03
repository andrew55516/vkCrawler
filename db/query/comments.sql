-- name: WriteDownComment :one
INSERT INTO comments (post_id,
                      owner,
                      thread_owner,
                      created_at)
VALUES ($1, $2, $3, $4) RETURNING *;

-- name: GetComment :one
SELECT * from comments
WHERE post_id = $1 AND owner = $2;

-- name: GetAllUnicCommentsOwners :many
SELECT  Distinct owner from comments;

-- name: GetAllCommentsByThreadOwner :many
SELECT owner, created_at from comments
WHERE thread_owner = $1;