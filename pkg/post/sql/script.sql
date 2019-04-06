CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE posts (
    uid UUID PRIMARY KEY,
    user_uid UUID NOT NULL,
    category_uid UUID NOT NULL,
    title VARCHAR(80) NOT NULL,
    url VARCHAR(80),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    modified_at TIMESTAMP WITH TIME ZONE NOT NULL
);

CREATE TABLE categories (
    uid UUID PRIMARY KEY,
    user_uid UUID NOT NULL,
    name VARCHAR(80) NOT NULL
);

ALTER TABLE posts
  ADD CONSTRAINT post_category_fk FOREIGN KEY (category_uid)
      REFERENCES categories(uid);
