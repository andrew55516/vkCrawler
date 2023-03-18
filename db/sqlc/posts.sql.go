// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.16.0
// source: posts.sql

package db

import (
	"context"
	"time"
)

const getAllPosts = `-- name: GetAllPosts :many
SELECT id, link, owner, likes, comments, reposts, created_at FROM posts
`

func (q *Queries) GetAllPosts(ctx context.Context) ([]Post, error) {
	rows, err := q.db.QueryContext(ctx, getAllPosts)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Post
	for rows.Next() {
		var i Post
		if err := rows.Scan(
			&i.ID,
			&i.Link,
			&i.Owner,
			&i.Likes,
			&i.Comments,
			&i.Reposts,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getLastPost = `-- name: GetLastPost :one
SELECT id, link, owner, likes, comments, reposts, created_at FROM posts
ORDER BY id DESC LIMIT 1
`

func (q *Queries) GetLastPost(ctx context.Context) (Post, error) {
	row := q.db.QueryRowContext(ctx, getLastPost)
	var i Post
	err := row.Scan(
		&i.ID,
		&i.Link,
		&i.Owner,
		&i.Likes,
		&i.Comments,
		&i.Reposts,
		&i.CreatedAt,
	)
	return i, err
}

const getPostByLink = `-- name: GetPostByLink :one
SELECT id, link, owner, likes, comments, reposts, created_at FROM posts
WHERE link = $1 LIMIT 1
`

func (q *Queries) GetPostByLink(ctx context.Context, link string) (Post, error) {
	row := q.db.QueryRowContext(ctx, getPostByLink, link)
	var i Post
	err := row.Scan(
		&i.ID,
		&i.Link,
		&i.Owner,
		&i.Likes,
		&i.Comments,
		&i.Reposts,
		&i.CreatedAt,
	)
	return i, err
}

const getUnknownPosts = `-- name: GetUnknownPosts :many
SELECT id, link, owner, likes, comments, reposts, created_at FROM posts
WHERE likes = -1
`

func (q *Queries) GetUnknownPosts(ctx context.Context) ([]Post, error) {
	rows, err := q.db.QueryContext(ctx, getUnknownPosts)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Post
	for rows.Next() {
		var i Post
		if err := rows.Scan(
			&i.ID,
			&i.Link,
			&i.Owner,
			&i.Likes,
			&i.Comments,
			&i.Reposts,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updatePost = `-- name: UpdatePost :one
UPDATE posts
SET owner = $1, likes = $2, comments = $3, reposts = $4
WHERE id = $5
RETURNING id, link, owner, likes, comments, reposts, created_at
`

type UpdatePostParams struct {
	Owner    string `json:"owner"`
	Likes    int64  `json:"likes"`
	Comments int64  `json:"comments"`
	Reposts  int64  `json:"reposts"`
	ID       int64  `json:"id"`
}

func (q *Queries) UpdatePost(ctx context.Context, arg UpdatePostParams) (Post, error) {
	row := q.db.QueryRowContext(ctx, updatePost,
		arg.Owner,
		arg.Likes,
		arg.Comments,
		arg.Reposts,
		arg.ID,
	)
	var i Post
	err := row.Scan(
		&i.ID,
		&i.Link,
		&i.Owner,
		&i.Likes,
		&i.Comments,
		&i.Reposts,
		&i.CreatedAt,
	)
	return i, err
}

const writeDownPost = `-- name: WriteDownPost :one
INSERT INTO posts (link,
                   owner, likes, comments, reposts, created_at)
VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, link, owner, likes, comments, reposts, created_at
`

type WriteDownPostParams struct {
	Link      string    `json:"link"`
	Owner     string    `json:"owner"`
	Likes     int64     `json:"likes"`
	Comments  int64     `json:"comments"`
	Reposts   int64     `json:"reposts"`
	CreatedAt time.Time `json:"created_at"`
}

func (q *Queries) WriteDownPost(ctx context.Context, arg WriteDownPostParams) (Post, error) {
	row := q.db.QueryRowContext(ctx, writeDownPost,
		arg.Link,
		arg.Owner,
		arg.Likes,
		arg.Comments,
		arg.Reposts,
		arg.CreatedAt,
	)
	var i Post
	err := row.Scan(
		&i.ID,
		&i.Link,
		&i.Owner,
		&i.Likes,
		&i.Comments,
		&i.Reposts,
		&i.CreatedAt,
	)
	return i, err
}