-- name: create :one
INSERT INTO
  users (
    username,
    email,
    password,
    bio,
    image,
    created_at,
    updated_at
  )
VALUES
  (?, ?, ?, ?, ?, ?, ?) RETURNING *;

-- name: getById :one
SELECT
  u.*,
  CAST(
    COALESCE(GROUP_CONCAT (f.follower_id, ','), '') AS TEXT
  ) AS followers
FROM
  users u
  LEFT JOIN followers f ON f.author_id = u.id
WHERE
  u.id = ?
GROUP BY
  u.id
LIMIT
  1;

-- name: getByEmail :one
SELECT
  u.*,
  CAST(
    COALESCE(GROUP_CONCAT (f.follower_id, ','), '') AS TEXT
  ) AS followers
FROM
  users u
  LEFT JOIN followers f ON f.author_id = u.id
WHERE
  u.email = ?
GROUP BY
  u.id
LIMIT
  1;

-- name: getByUsername :one
SELECT
  u.*,
  CAST(
    COALESCE(GROUP_CONCAT (f.follower_id, ','), '') AS TEXT
  ) AS followers
FROM
  users u
  LEFT JOIN followers f ON f.author_id = u.id
WHERE
  u.username = ?
GROUP BY
  u.id
LIMIT
  1;

-- name: getFollowers :many
SELECT
  follower_id
FROM
  followers
WHERE
  author_id = ?;

-- name: getFollowing :many
SELECT
  f.author_id
FROM
  followers f
WHERE
  f.follower_id = sqlc.arg (author_id);

-- name: update :one
UPDATE users
SET
  username = ?,
  email = ?,
  password = ?,
  bio = ?,
  image = ?,
  created_at = ?,
  updated_at = ?
WHERE
  id = ? RETURNING *;

-- name: follow :execrows
INSERT INTO
  followers (author_id, follower_id, created_at)
VALUES
  (?, ?, ?);

-- name: unfollow :execrows
DELETE FROM followers
WHERE
  author_id = ?
  AND follower_id = ?;
