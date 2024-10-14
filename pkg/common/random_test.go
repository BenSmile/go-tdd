package common

import "testing"

func TestRandomString(t *testing.T) {
	s := RandomString(32)
	if len(s) != 32 {
		t.Errorf("expected RandomString to procude a string of len 32")
	}
}
