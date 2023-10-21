-- only need create, read, delete
-- name: create :one
INSERT INTO
  comments (body, article_id, author_id, created_at)
VALUES
  (?, ?, ?, ?) RETURNING *;

-- name: getById :one
SELECT
  *
FROM
  comments
WHERE
  id = ?;

-- name: getByArticleId :many
SELECT
  *
FROM
  comments
WHERE
  article_id = ?;

-- name: getByAuthorId :many
SELECT
  *
FROM
  comments
WHERE
  author_id = ?;

-- name: delete :exec
DELETE FROM comments
WHERE
  id = ?
  AND author_id = ?;
