package postgres

import (
	"testing"

	"github.com/bensmile/go-api-tdd/pkg/domain"
)

var (
	oldSqlDeleteAllPosts     = sqlDeleteAllPosts
	oldSqlCreatePost         = sqlCreatePost
	oldSqlSelectPostById     = sqlSelectPostById
	oldSqlSelectPostByUserId = sqlSelectPostByUserId
)

func TestCreatePost(t *testing.T) {

	pStore := NewPostgresStore(testDB)

	err := pStore.DeleteAllPosts()

	if err != nil {
		t.Fatal(err)
	}

	err = pStore.DeleteAllUsers()

	if err != nil {
		t.Fatal(err)
	}

	user := &domain.User{
		Email:    "test@test.com",
		Password: "password",
		Name:     "JOhn Doe",
	}

	createdUser, err := pStore.CreateUser(user)
	if err != nil {
		t.Fatal(err)
	}

	post := &domain.Post{
		UserId: createdUser.Id,
		Title:  "test",
		Body:   "test body",
	}

	createdPost, err := pStore.CreatePost(post)

	if err != nil {
		t.Fatal(err)
	}

	if createdPost.Id == 0 {
		t.Errorf("want id not to be zero")
	}

	if post.Title != createdPost.Title {
		t.Errorf("expected title %q; got %q", post.Title, createdPost.Title)
	}

	if post.UserId != user.Id {
		t.Errorf("expected user id %q; got %q", user.Id, post.UserId)
	}

	err = pStore.DeleteAllPosts()

	if err != nil {
		t.Fatal(err)
	}

	sqlCreatePost = "invalid query"

	_, err = pStore.CreatePost(post)

	if err == nil {
		t.Errorf("expected error when CreatePost %v", err)
	}

	sqlCreatePost = oldSqlCreatePost

}

func TestDeleteAllPost(t *testing.T) {

	pStore := NewPostgresStore(testDB)

	_ = pStore.DeleteAllUsers()

	user := &domain.User{
		Email:    "test@test.com",
		Password: "password",
		Name:     "JOhn Doe",
	}

	createdUser, err := pStore.CreateUser(user)

	if err != nil {
		t.Fatal(err)
	}

	post := &domain.Post{
		UserId: createdUser.Id,
		Title:  "test",
		Body:   "test body",
	}

	_, err = pStore.CreatePost(post)

	if err != nil {
		t.Errorf("expected no error; got %q", err)
	}

	posts, err := pStore.FindPostsByUser(createdUser.Id)

	if err != nil {
		t.Errorf("expected no error; got %v", err)
	}
	if len(posts) == 0 {
		t.Errorf("expected posts for this user before deleting all posts")
	}

	err = pStore.DeleteAllPosts()

	if err != nil {
		t.Errorf("expecting no error for DeleteAllPosts; got %v", err)
	}

	posts, err = pStore.FindPostsByUser(createdUser.Id)

	if err != nil {
		t.Errorf("expected no error; got %v", err)
	}

	if len(posts) != 0 {
		t.Errorf("expected not posts for this user after deleting all posts")
	}

	sqlDeleteAllPosts = "invalid"

	err = pStore.DeleteAllPosts()
	if err == nil {
		t.Errorf("want error; got nil for invalid DeleteAllPosts sql")
	}

	sqlDeleteAllPosts = oldSqlDeleteAllPosts

	_ = pStore.DeleteAllUsers()
}

func TestFindPostById(t *testing.T) {
	pStore := NewPostgresStore(testDB)

	if err := pStore.DeleteAllPosts(); err != nil {
		t.Fatal(err)
	}

	if err := pStore.DeleteAllUsers(); err != nil {
		t.Fatal(err)
	}

	user := &domain.User{
		Email:    "test@test.com",
		Password: "password",
		Name:     "JOhn Doe",
	}

	createdUser, err := pStore.CreateUser(user)

	if err != nil {
		t.Fatal(err)
	}

	post := &domain.Post{
		UserId: createdUser.Id,
		Title:  "test",
		Body:   "test body",
	}

	newPost, err := pStore.CreatePost(post)

	if err != nil {
		t.Errorf("expected no error; got %q", err)
	}

	postById, err := pStore.FindPostById(newPost.Id)

	if err != nil {
		t.Errorf("expected no error; got %q", err)
	}

	if postById.Id == 0 {
		t.Errorf("want id not to be zero")
	}

	sqlSelectPostById = "invalid sql"

	_, err = pStore.FindPostById(newPost.Id)

	if err == nil {
		t.Errorf("expected error; got nil for invalid sql")
	}

	sqlSelectPostById = `SELECT id FROM posts WHERE id = $1`

	_, err = pStore.FindPostById(newPost.Id)

	if err == nil {
		t.Errorf("expected error; scan more than returned fields")
	}

	sqlSelectPostById = oldSqlSelectPostById

	_, err = pStore.FindPostById(-1)

	if err == nil {
		t.Errorf("expected error; post not found for the given id")
	}
}

func TestFindPostsByUser(t *testing.T) {
	pStore := NewPostgresStore(testDB)

	if err := pStore.DeleteAllPosts(); err != nil {
		t.Fatal(err)
	}

	if err := pStore.DeleteAllUsers(); err != nil {
		t.Fatal(err)
	}

	user := &domain.User{
		Email:    "test@test.com",
		Password: "password",
		Name:     "JOhn Doe",
	}

	createdUser, err := pStore.CreateUser(user)

	if err != nil {
		t.Fatal(err)
	}

	post := &domain.Post{
		UserId: createdUser.Id,
		Title:  "test",
		Body:   "test body",
	}

	newPost, err := pStore.CreatePost(post)

	posts, err := pStore.FindPostsByUser(newPost.UserId)

	if err != nil {
		t.Errorf("expected no error; got %v", err)
	}

	if len(posts) == 0 {
		t.Errorf("expected posts for this user")
	}

	sqlSelectPostByUserId = "invalid sql"

	_, err = pStore.FindPostsByUser(createdUser.Id)

	if err == nil {
		t.Errorf("expected error; got nil for invalid sql for FindPostsByUser")
	}

	sqlSelectPostByUserId = `SELECT id FROM posts WHERE user_id = $1`

	_, err = pStore.FindPostsByUser(createdUser.Id)

	if err == nil {
		t.Errorf("expected error; scan more than returned fields for FindPostsByUser")
	}

	sqlSelectPostByUserId = oldSqlSelectPostByUserId
}
