package auth

import (
	"testing"
	"time"

	"employee-directory-api/internal/domain"

	"github.com/google/uuid"
)

func TestHashAndCheckPassword(t *testing.T) {
	pw := "secretPassword123!"
	hash, err := HashPassword(pw)
	if err != nil {
		t.Fatalf("HashPassword failed: %v", err)
	}
	if !CheckPassword(hash, pw) {
		t.Fatalf("expected valid password check")
	}
	if CheckPassword(hash, "wrongpassword") {
		t.Fatalf("expected invalid password check to fail")
	}
}

func TestIssueAndParseToken(t *testing.T) {
	svc, err := NewService("test-secret-key-12345", "1h")
	if err != nil {
		t.Fatalf("NewService failed: %v", err)
	}

	accID := uuid.New()
	personID := uuid.New()
	orgID := uuid.New()
	acc := &domain.PersonAccount{
		ID:       accID,
		Username: "tester",
		PersonID: personID,
		Role:     "manager",
	}

	tok, err := svc.IssueToken(acc, orgID)
	if err != nil {
		t.Fatalf("IssueToken failed: %v", err)
	}

	claims, err := svc.ParseToken(tok)
	if err != nil {
		t.Fatalf("ParseToken failed: %v", err)
	}
	if claims.Username != "tester" || claims.PersonID != personID || claims.OrgID != orgID || claims.Role != "manager" {
		t.Fatalf("claims mismatch: %+v", claims)
	}
	if claims.Subject != accID.String() {
		t.Fatalf("expected subject %s, got %s", accID.String(), claims.Subject)
	}
}

func TestParseToken_InvalidSecret(t *testing.T) {
	svc1, _ := NewService("secret-one", "1h")
	svc2, _ := NewService("secret-two", "1h")

	acc := &domain.PersonAccount{
		ID:       uuid.New(),
		Username: "alice",
		PersonID: uuid.New(),
		Role:     "user",
	}

	tok, err := svc1.IssueToken(acc, uuid.New())
	if err != nil {
		t.Fatalf("IssueToken failed: %v", err)
	}

	_, err = svc2.ParseToken(tok)
	if err == nil {
		t.Fatalf("expected error parsing token with different secret, got nil")
	}
}

func TestParseToken_Expired(t *testing.T) {
	svc, _ := NewService("secret-key", "-1h")

	acc := &domain.PersonAccount{
		ID:       uuid.New(),
		Username: "bob",
		PersonID: uuid.New(),
		Role:     "admin",
	}

	tok, err := svc.IssueToken(acc, uuid.New())
	if err != nil {
		t.Fatalf("IssueToken failed: %v", err)
	}

	time.Sleep(10 * time.Millisecond)
	_, err = svc.ParseToken(tok)
	if err == nil {
		t.Fatalf("expected error parsing expired token, got nil")
	}
}
