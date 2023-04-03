-- name: WriteDownLike :one
INSERT INTO likes (post_id,
                      owner)
VALUES ($1, $2) RETURNING *;

-- name: GetLike :one
SELECT * from likes
WHERE post_id = $1 AND owner = $2;

-- name: GetAllUnicLikesOwners :many
SELECT  Distinct owner from likes;

-- name: GetAllLikesByPostID :many
SELECT owner from likes
WHERE post_id = $1;