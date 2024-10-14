-- Create the blog_test database
CREATE DATABASE blog_test;

-- Switch to the 'blog' database and create the 'users' and 'posts' tables
\connect blog;

CREATE TABLE users(
                      id bigserial PRIMARY KEY,
                      name varchar(255) NOT NULL,
                      email varchar(255) NOT NULL UNIQUE,
                      password varchar(255) NOT NULL
);

CREATE TABLE posts(
                      id bigserial PRIMARY KEY,
                      title TEXT NOT NULL UNIQUE,
                      body TEXT NOT NULL,
                      user_id bigint NOT NULL,
                      created_at TIMESTAMPTZ NOT NULL DEFAULT (now()),
                      FOREIGN KEY (user_id) REFERENCES users (id)
);

-- Switch to the 'blog_test' database and create the same schema for testing purposes
\connect blog_test;

CREATE TABLE users(
                      id bigserial PRIMARY KEY,
                      name varchar(255) NOT NULL,
                      email varchar(255) NOT NULL UNIQUE,
                      password varchar(255) NOT NULL
);

CREATE TABLE posts(
                      id bigserial PRIMARY KEY,
                      title TEXT NOT NULL UNIQUE,
                      body TEXT NOT NULL,
                      user_id bigint NOT NULL,
                      created_at TIMESTAMPTZ NOT NULL DEFAULT (now()),
                      FOREIGN KEY (user_id) REFERENCES users (id)
);