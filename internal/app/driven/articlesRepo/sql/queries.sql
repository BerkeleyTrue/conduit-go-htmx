-- name: create :one
INSERT INTO
  articles (
    slug,
    title,
    description,
    body,
    author_id,
    created_at,
    updated_at
  )
VALUES
  (?, ?, ?, ?, ?, ?, ?) RETURNING *;

-- name: createTag :execrows
INSERT INTO
  tags (tag)
VALUES
  (?) ON CONFLICT DO NOTHING;

-- name: createArticleTag :execrows
INSERT INTO
  article_tags (article_id, tag_id)
VALUES
  (
    ?,
    (
      SELECT
        id
      FROM
        tags
      WHERE
        tag = ?
    )
  );

-- name: getBySlug :one
SELECT
  sqlc.embed(a),
  GROUP_CONCAT(t.tag, ',') AS tags
FROM
  articles a
  LEFT JOIN article_tags at ON a.id = at.article_id
  LEFT JOIN tags t ON at.tag_id = t.id
WHERE
  slug = ?
LIMIT
  1;

-- name: getById :one
SELECT
  sqlc.embed(a),
  GROUP_CONCAT(t.tag, ',') AS tags
FROM
  articles a
  LEFT JOIN article_tags at ON a.id = at.article_id
  LEFT JOIN tags t ON at.tag_id = t.id
WHERE
  articles.id = ?
LIMIT
  1;

-- name: list :many
SELECT
  sqlc.embed(a),
  GROUP_CONCAT(t.tag, ',') AS tags
FROM
  articles a
  LEFT JOIN article_tags at ON a.id = at.article_id
  LEFT JOIN tags t ON at.tag_id = t.id
WHERE
  (cast(sqlc.narg(tag) as text) IS NULL OR a.id IN (
    SELECT
      at2.article_id
    FROM
      article_tags at2
      LEFT JOIN tags t2 ON t2.id = at2.tag_id
    WHERE
    t2.tag = sqlc.narg(tag)
    )
  )
  AND (cast(sqlc.narg(author_id) as integer) IS NULL OR a.author_id = sqlc.narg(author_id))
  AND (cast(sqlc.narg(favorited) as integer) IS NULL OR a.id IN (
    SELECT
      f.article_id
    FROM
      favorites f
    WHERE
      f.user_id = sqlc.narg(favorited)
    )
  )
GROUP BY
  a.id
ORDER BY
  a.created_at DESC
LIMIT
  sqlc.arg(limit)
OFFSET
  sqlc.arg(offset);

-- name: feed :many
SELECT
  sqlc.embed(a),
  GROUP_CONCAT(t.tag, ',') AS tags
FROM
  articles a
  LEFT JOIN article_tags at ON a.id = at.article_id
  LEFT JOIN tags t ON at.tag_id = t.id
WHERE
  a.author_id in (
    SELECT
      f.author_id
    FROM
      followers f
    WHERE
      f.follower_id = sqlc.arg(followed)
  )
GROUP BY
  a.id
ORDER BY
  a.created_at DESC
LIMIT
  sqlc.arg(limit)
OFFSET
  sqlc.arg(offset);

-- name: getPopularTags :many
SELECT
  tags.tag
FROM
  tags
  LEFT JOIN article_tags ON tags.id = article_tags.tag_id
GROUP BY
  tags.id,
  tags.tag
ORDER BY
  COUNT(article_tags.tag_id) DESC
LIMIT
  10;

-- name: update :one
UPDATE articles
SET
  title = ?,
  slug = ?,
  description = ?,
  body = ?,
  updated_at = ?
WHERE
  articles.id = ? RETURNING *;

-- name: favorite :execrows
INSERT INTO
  favorites (user_id, article_id)
VALUES
  (?, ?);

-- name: unfavorite :execrows
DELETE FROM favorites
WHERE
  user_id = ?
  AND article_id = ?;

-- name: delete :execrows
DELETE FROM articles
WHERE
  slug = ?;

-- name: getNumOfFavorites :one
SELECT
  COUNT(*) as count
FROM
  favorites f
WHERE
  article_id = ?;

-- name: isFavoritedByUser :one
SELECT
  COUNT(*) as count
FROM
  favorites f
WHERE
  f.user_id = ?
  AND f.article_id = ?;
