package postgres

import (
	"errors"
	"testing"

	"github.com/bensmile/go-api-tdd/pkg/domain"
)

var (
	oldSqlCreateUser        = sqlCreateUser
	oldSqlDeleteUserById    = sqlDeleteUserById
	oldSqlSelectUserByEmail = sqlSelectUserByEmail
	oldSqlSelectUserById    = sqlSelectUserById
	oldSqlDeleteAllUsers    = sqlDeleteAllUserss
)

func TestCreateUser(t *testing.T) {

	pStore := NewPostgresStore(testDB)

	if err := pStore.DeleteAllPosts(); err != nil {
		t.Fatal(err)
	}

	if err := pStore.DeleteAllUsers(); err != nil {
		t.Fatal(err)
	}

	oldPassword := "password"

	user := &domain.User{
		Email:    "test@test.com",
		Password: oldPassword,
		Name:     "JOhn Doe",
	}

	createdUser, err := pStore.CreateUser(user)

	if err != nil {
		t.Fatal(err)
	}

	if createdUser.Id == 0 {
		t.Errorf("want id not to be zero")
	}

	if user.Name != createdUser.Name {
		t.Errorf("expected %q; got %q", user.Name, createdUser.Name)
	}

	if createdUser.Password == oldPassword {
		t.Error("password was not hashed")
	}

	sqlCreateUser = "invalid query"

	_, err = pStore.CreateUser(user)
	if err == nil {
		t.Errorf("expected error not to be nil for invalid CreateUser sql")
	}

	sqlCreateUser = oldSqlCreateUser

	err = pStore.DeleteUserById(createdUser.Id)

	if err != nil {
		t.Errorf("expected nil error during DeleteUserById; got %q", err)
	}

	sqlDeleteUserById = "invalid query"

	err = pStore.DeleteUserById(createdUser.Id)
	if err == nil {
		t.Errorf("expected error not to be nil for invalid DeleteUserById sql")
	}

	sqlDeleteUserById = oldSqlDeleteUserById

}

func TestFindUserByEmail(t *testing.T) {
	pStore := NewPostgresStore(testDB)
	user := &domain.User{
		Email:    "test@test.com",
		Password: "password",
		Name:     "JOhn Doe",
	}

	createdUser, err := pStore.CreateUser(user)
	if err != nil {
		t.Fatal(err)
	}

	userByEmail, err := pStore.FindUserByEmail(createdUser.Email)

	if err != nil {
		t.Errorf("expected no error; got %q", err)
	}

	if userByEmail.Email != createdUser.Email {
		t.Errorf("expect %q; got %q", createdUser.Email, userByEmail.Email)
	}

	_, err = pStore.FindUserByEmail("invalid email")
	if err == nil {
		t.Errorf("want error; got nil for invalid email")
	}

	if err != nil && !errors.Is(err, domain.ErrUserNotFound) {
		t.Errorf("want domain.ErrUserNotFound error; got %q", err)
	}

	sqlSelectUserByEmail = "invalid"

	_, err = pStore.FindUserByEmail(createdUser.Email)
	if err == nil {
		t.Errorf("want error; got nil for invalid FindUserByEmail sql")
	}

	sqlSelectUserByEmail = oldSqlSelectUserByEmail

	_ = pStore.DeleteUserById(createdUser.Id)

}

func TestFindUserById(t *testing.T) {
	pStore := NewPostgresStore(testDB)
	user := &domain.User{
		Email:    "test@test.com",
		Password: "password",
		Name:     "JOhn Doe",
	}

	createdUser, err := pStore.CreateUser(user)
	if err != nil {
		t.Fatal(err)
	}

	userById, err := pStore.FindUserById(createdUser.Id)

	if err != nil {
		t.Errorf("expected no error; got %q", err)
	}

	if userById.Id != createdUser.Id {
		t.Errorf("expect %q; got %q", createdUser.Email, userById.Email)
	}

	_, err = pStore.FindUserById(-1)
	if err == nil {
		t.Errorf("want error; got nil for invalid id")
	}

	if err != nil && !errors.Is(err, domain.ErrUserNotFound) {
		t.Errorf("want domain.ErrUserNotFound error; got %q", err)
	}

	sqlSelectUserById = "invalid"

	_, err = pStore.FindUserById(createdUser.Id)
	if err == nil {
		t.Errorf("want error; got nil for invalid FindUserById sql")
	}

	sqlSelectUserById = oldSqlSelectUserById

	_ = pStore.DeleteUserById(createdUser.Id)

}

func TestDeleteAllUsers(t *testing.T) {
	pStore := NewPostgresStore(testDB)
	user := &domain.User{
		Email:    "test@test.com",
		Password: "password",
		Name:     "JOhn Doe",
	}

	_, err := pStore.CreateUser(user)
	if err != nil {
		t.Fatal(err)
	}

	err = pStore.DeleteAllUsers()

	if err != nil {
		t.Errorf("expected no error; got %q", err)
	}

	sqlDeleteAllUserss = "invalid"

	err = pStore.DeleteAllUsers()
	if err == nil {
		t.Errorf("want error; got nil for invalid DeleteAllUsers sql")
	}

	sqlDeleteAllUserss = oldSqlDeleteAllUsers

}
