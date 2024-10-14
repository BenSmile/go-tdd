package postgres

import (
	"database/sql"
	"errors"

	"github.com/bensmile/go-api-tdd/pkg/domain"
)

var (
	sqlCreatePost = `INSERT INTO posts (user_id, title, body)
								 VALUES($1, $2, $3) RETURNING *`
	sqlSelectPostByUserId = `SELECT * FROM posts WHERE user_id = $1`
	sqlSelectPostById     = `SELECT * FROM posts WHERE id = $1`
	sqlDeleteAllPosts     = `DELETE FROM posts`
)

func (q *postgresStore) CreatePost(post *domain.Post) (*domain.Post, error) {
	err := q.db.QueryRow(sqlCreatePost, post.UserId, post.Title, post.Body).
		Scan(&post.Id, &post.Title, &post.Body, &post.UserId, &post.CreatedAt)
	if err != nil {
		return nil, err
	}
	return post, nil
}

func (q *postgresStore) FindPostsByUser(userId int64) ([]domain.Post, error) {
	var posts []domain.Post
	rows, err := q.db.Query(sqlSelectPostByUserId, userId)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var post domain.Post
		err := rows.Scan(&post.Id, &post.Title, &post.Body, &post.UserId, &post.CreatedAt)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	return posts, nil
}

func (q *postgresStore) FindPostById(id int64) (*domain.Post, error) {
	var post domain.Post
	err := q.db.QueryRow(sqlSelectPostById, id).
		Scan(&post.Id, &post.Title, &post.Body, &post.UserId, &post.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrPostNotFound
		}
		return nil, err
	}
	return &post, nil
}

func (q *postgresStore) DeleteAllPosts() error {
	_, err := q.db.Exec(sqlDeleteAllPosts)
	if err != nil {
		return err
	}
	return nil
}
