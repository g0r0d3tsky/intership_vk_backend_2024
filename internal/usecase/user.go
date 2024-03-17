package usecase

import (
	"cinema_service/internal/domain"
	"context"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"time"
)

const (
	salt       = "some-salt"
	signingKey = "some-signing-key"
	tokenTTL   = 12 * time.Hour
)

type UserInfo struct {
	UserID uuid.UUID `json:"user_id"`
	Role   string    `json:"role"`
}

type tokenClaims struct {
	jwt.RegisteredClaims
	userClaims *UserInfo
}

//go:generate mockgen -source=user.go -destination=mocks/userMock.go

type UserRepo interface {
	GetUser(ctx context.Context, login string, password string) (*domain.User, error)
	//CreateUser(ctx context.Context, u *domain.User) error
	//DeleteUser(ctx context.Context, userID uuid.UUID) error
}

type UserService struct {
	repo UserRepo
}

func NewUserService(repo UserRepo) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) GetUser(ctx context.Context, login string, password string) (*domain.User, error) {
	user, err := s.repo.GetUser(ctx, login, password)
	if err != nil {
		return nil, fmt.Errorf("getting user: %w", err)
	}
	return user, nil
}

// TODO: validation?
// func (s *UserService) CreateUser(ctx context.Context, u *domain.User) error {
// 	err := s.repo.CreateUser(ctx, u)
// 	if err != nil {
// 		return fmt.Errorf("create user: %w", err)
// 	}

// 	return nil
// }

// func (s *UserService) DeleteUser(ctx context.Context, userID uuid.UUID) error {
// 	err := s.repo.DeleteUser(ctx, userID)
// 	if err != nil {
// 		return fmt.Errorf("delete user: %w", err)
// 	}

// 	return nil
// }

func (s *UserService) GenerateToken(ctx context.Context, login string, password string) (string, error) {
	user, err := s.GetUser(ctx, login, password)
	if err != nil {
		return "", fmt.Errorf("get user: %w", err)
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		&UserInfo{user.ID,
			user.Role,
		},
	})

	signedToken, err := token.SignedString([]byte(signingKey))
	if err != nil {
		return "", fmt.Errorf("sign token: %w", err)
	}

	return signedToken, nil
}

func (s *UserService) ParseToken(token string) (*UserInfo, error) {
	t, err := jwt.ParseWithClaims(token, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(signingKey), nil
	})
	if err != nil {
		return nil, fmt.Errorf("parse token: %w", err)
	}

	claims, ok := t.Claims.(*tokenClaims)
	if !ok {
		return nil, errors.New("token claims are not of type *tokenClaims")
	}

	return claims.userClaims, nil
}
