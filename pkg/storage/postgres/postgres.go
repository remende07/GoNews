package postgres

import (
	"GoNews/pkg/storage"
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
)

type Store struct {
	db *pgxpool.Pool
}

func New(connstr string) (*Store, error) {
	db, err := pgxpool.Connect(context.Background(), connstr)
	if err != nil {
		return nil, err
	}

	s := Store{
		db: db,
	}

	return &s, nil
}

func (s *Store) Posts() ([]storage.Post, error) {
	rows, err := s.db.Query(context.Background(), `
	SELECT
	posts.id,
    posts.title,
    posts.content,
	authors.id,
	authors.name,
    posts.created_at
	FROM posts, authors
	`)

	if err != nil {
		return nil, err
	}

	var posts []storage.Post

	for rows.Next() {
		var p storage.Post
		err = rows.Scan(
			&p.ID,
			&p.Title,
			&p.Content,
			&p.AuthorID,
			&p.AuthorName,
			&p.CreatedAt,
		)

		if err != nil {
			return nil, err
		}

		posts = append(posts, p)
	}
	return posts, rows.Err()
}

func (s *Store) AddPost(post storage.Post) error {
	_, err := s.db.Exec(context.Background(), `
		INSERT INTO posts(author_id, title, content)
		VALUES
		($1, $2, $3);
	`, post.AuthorID, post.Title, post.Content)

	if err != nil {
		return err
	}
	return nil
}

func (s *Store) UpdatePost(post storage.Post) error {
	_, err := s.db.Exec(context.Background(), `
		UPDATE posts
		SET title = $1, content = $2
		WHERE id = $3;
	`, post.Title, post.Content, post.ID)

	if err != nil {
		return err
	}
	return nil
}

func (s *Store) DeletePost(post storage.Post) error {
	_, err := s.db.Exec(context.Background(), `
		DELETE FROM posts WHERE id = $1;
	`, post.ID)

	if err != nil {
		return err
	}
	return nil
}
