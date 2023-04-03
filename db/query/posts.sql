-- name: WriteDownPost :one
INSERT INTO posts (link,
                   owner, likes, comments, reposts, created_at)
VALUES ($1, $2, $3, $4, $5, $6) RETURNING *;

-- name: GetPostByLink :one
SELECT * FROM posts
WHERE link = $1 LIMIT 1;

-- name: GetLastPost :one
SELECT * FROM posts
ORDER BY id DESC LIMIT 1;

-- name: GetAllPosts :many
SELECT * FROM posts;

-- name: UpdatePost :one
UPDATE posts
SET owner = $1, likes = $2, comments = $3, reposts = $4
WHERE id = $5
RETURNING *;

-- name: GetUnknownPosts :many
SELECT * FROM posts
WHERE likes = -1;

-- name: GetAllUnicPostsOwners :many
SELECT  Distinct owner from posts;




