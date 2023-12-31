CREATE TABLE
  IF NOT EXISTS users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username TEXT NOT NULL UNIQUE,
    email TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL,
    bio TEXT,
    image TEXT,
    created_at TEXT NOT NULL,
    updated_at TEXT
  );

CREATE TABLE
  IF NOT EXISTS followers (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    author_id INTEGER NOT NULL,
    follower_id INTEGER NOT NULL,
    created_at TEXT NOT NULL,
    UNIQUE (author_id, follower_id),
    FOREIGN KEY (author_id) REFERENCES users (id),
    FOREIGN KEY (follower_id) REFERENCES users (id)
  );
