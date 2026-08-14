package domain

import "testing"

func TestHashPasswordRoundTrip(t *testing.T) {
	hash, err := HashPassword("correct-horse-battery-staple")
	if err != nil {
		t.Fatalf("HashPassword returned error: %v", err)
	}
	if hash == "" {
		t.Fatal("HashPassword returned an empty hash")
	}
	if hash == "correct-horse-battery-staple" {
		t.Fatal("HashPassword returned the plaintext password unchanged")
	}
	if !CheckPasswordHash("correct-horse-battery-staple", hash) {
		t.Fatal("CheckPasswordHash rejected the correct password")
	}
}

func TestCheckPasswordHashWrongPassword(t *testing.T) {
	hash, err := HashPassword("correct-horse-battery-staple")
	if err != nil {
		t.Fatalf("HashPassword returned error: %v", err)
	}
	if CheckPasswordHash("wrong-password", hash) {
		t.Fatal("CheckPasswordHash accepted an incorrect password")
	}
}

func TestCheckPasswordHashEmptyPassword(t *testing.T) {
	hash, err := HashPassword("correct-horse-battery-staple")
	if err != nil {
		t.Fatalf("HashPassword returned error: %v", err)
	}
	if CheckPasswordHash("", hash) {
		t.Fatal("CheckPasswordHash accepted an empty password")
	}
}
