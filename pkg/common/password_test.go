package common

import "testing"

func TestPasswordHash(t *testing.T) {

	password := "password"
	hashed, err := HashPassword(password)

	if err != nil {
		t.Fatal(err)
	}

	if len(hashed) == 0 {
		t.Errorf("want hash; got %q", hashed)
	}

	if hashed == password {
		t.Errorf("password was not hashed; got %q", hashed)
	}
}

func TestPasswordHash_WithError(t *testing.T) {

	longPass := make([]byte, 73)

	_, err := HashPassword(string(longPass))

	if err == nil {
		t.Errorf("expected an error; got %q", err)
	}

}

func TestCheckPassword(t *testing.T) {
	password := "password"

	hashed, err := HashPassword(password)

	if err != nil {
		t.Fatal(err)
	}

	err = CheckPassword(password, hashed)

	if err != nil {
		t.Errorf("password verification failed; got %q", err)
	}
}
