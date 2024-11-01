DROP TABLE IF EXISTS posts, authors;

CREATE TABLE authors (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL
);

CREATE TABLE posts (
    id SERIAL PRIMARY KEY,
    author_id INTEGER REFERENCES authors(id) NOT NULL,
    title TEXT NOT NULL,
    content TEXT NOT NULL,
    created_at BIGINT DEFAULT extract(epoch from now()),
);

INSERT INTO authors (id, name) VALUES (0, 'Дмитрий');
INSERT INTO authors (id, name) VALUES (1, 'Мария');
INSERT INTO posts (id, author_id, title, content, created_at) VALUES (0, 0, 'Статья', 'Содержание статьи', 0);