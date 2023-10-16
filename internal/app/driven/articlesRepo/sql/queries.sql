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
  *
FROM
  articles
WHERE
  slug = ?
LIMIT
  1;

-- name: getById :one
SELECT
  *
FROM
  articles
WHERE
  id = ?
LIMIT
  1;

-- name: list :many
SELECT
  *
FROM
  articles
  LEFT JOIN article_tags ON articles.id = article_tags.article_id
  LEFT JOIN tags ON article_tags.tag_id = tags.id
WHERE
  author_id = coalesce(sqlc.narg(author_id), author_id)
GROUP BY
  articles.id
ORDER BY
  articles.created_at DESC
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
  description = ?,
  body = ?,
  updated_at = ?
WHERE
  slug = ? RETURNING *;

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
