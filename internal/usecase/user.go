package usecase

import (
	"cinema_service/internal/domain"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const (
	//salt       = "some-salt"
	signingKey = "some-signing-key"
	tokenTTL   = 12 * time.Hour
)

type UserInfo struct {
	UserID uuid.UUID `json:"user_id"`
	Role   string    `json:"role"`
}

type tokenClaims struct {
	jwt.RegisteredClaims
	UserClaims UserInfo `json:"userClaims"`
}

//go:generate mockgen -source=user.go -destination=mocks/userMock.go

type UserRepo interface {
	GetUser(ctx context.Context, login string, password string) (*domain.User, error)
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

func (s *UserService) GenerateToken(ctx context.Context, login string, password string) (string, error) {
	user, err := s.GetUser(ctx, login, password)
	if err != nil {
		return "", fmt.Errorf("get user: %w", err)
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		UserClaims: UserInfo{
			UserID:   user.ID,
			Role: user.Role,
		},
	})

	signedToken, err := token.SignedString([]byte(signingKey))
	if err != nil {
		return "", fmt.Errorf("sign token: %w", err)
	}

	return signedToken, nil
}

func (s *UserService) ParseToken(accessToken string) (*UserInfo, error) {
	token, err := jwt.ParseWithClaims(accessToken, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}

		return []byte(signingKey), nil
	})
	if err != nil {
		return nil, err
	}
	
	if claims, ok := token.Claims.(*tokenClaims); ok && token.Valid {
		return &claims.UserClaims, nil
	}

	return nil, errors.New("invalid token")
}