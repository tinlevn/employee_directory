package auth

import (
	"errors"
	"time"

	"employee-directory-api/internal/domain"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Claims struct {
	Username string    `json:"username"`
	PersonID uuid.UUID `json:"person_id"`
	OrgID    uuid.UUID `json:"org_id"`
	Role     string    `json:"role"`
	jwt.RegisteredClaims
}

type Service struct {
	secret []byte
	ttl    time.Duration
}

func NewService(secret, ttl string) (*Service, error) {
	d, err := time.ParseDuration(ttl)
	if err != nil {
		d = 24 * time.Hour
	}
	return &Service{secret: []byte(secret), ttl: d}, nil
}

func (s *Service) IssueToken(account *domain.PersonAccount, orgID uuid.UUID) (string, error) {
	now := time.Now()
	claims := Claims{
		Username: account.Username,
		PersonID: account.PersonID,
		OrgID:    orgID,
		Role:     account.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   account.ID.String(),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(s.ttl)),
			Issuer:    "employee-directory-api",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.secret)
}

func (s *Service) ParseToken(tokenString string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return s.secret, nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, errors.New("invalid token")
	}
	return claims, nil
}

func HashPassword(password string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func CheckPassword(hash, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}
