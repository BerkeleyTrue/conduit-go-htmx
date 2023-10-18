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
  *
FROM
  users
WHERE
  id = ?
LIMIT
  1;

-- name: getByEmail :one
SELECT
  *
FROM
  users
WHERE
  email = ?
LIMIT
  1;

-- name: getByUsername :one
SELECT
  *
FROM
  users
WHERE
  username = ?
LIMIT
  1;

-- name: getFollowers :many
SELECT
  follower_id
FROM
  followers
WHERE
  user_id = ?;


-- name: getFollowing :many
SELECT
  f.user_id
FROM
  followers f
WHERE
  f.follower_id = sqlc.arg(user_id);

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
  followers (user_id, follower_id)
VALUES
  (?, ?);

-- name: unfollow :execrows
DELETE FROM followers
WHERE
  user_id = ?
  AND follower_id = ?;
