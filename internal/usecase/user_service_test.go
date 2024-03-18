package usecase

import (
	"cinema_service/internal/domain"
	mock_repo "cinema_service/internal/usecase/mocks"
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestGetUser(t *testing.T) {
	tests := []struct {
		name          string
		login         string
		password      string
		mockUser      *domain.User
		mockError     error
		expectedUser  *domain.User
		expectedError error
	}{
		{
			name:          "User found",
			login:         "existinguser",
			password:      "correctpassword",
			mockUser:      &domain.User{ID: uuid.UUID{00000000 - 0000 - 0000 - 0000 - 000000000000}, Role: "user"},
			mockError:     nil,
			expectedUser:  &domain.User{ID: uuid.UUID{00000000 - 0000 - 0000 - 0000 - 000000000000}, Role: "user"},
			expectedError: nil,
		},
		{
			name:          "User not found",
			login:         "nonexistentuser",
			password:      "password",
			mockUser:      nil,
			mockError:     errors.New("user not found"),
			expectedUser:  nil,
			expectedError: errors.New("getting user: user not found"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			mockUserRepo := mock_repo.NewMockUserRepo(c)
			service := NewUserService(mockUserRepo)

			mockUserRepo.EXPECT().GetUser(gomock.Any(), test.login, test.password).Return(test.mockUser, test.mockError)

			user, err := service.GetUser(context.Background(), test.login, test.password)
			assert.Equal(t, test.expectedUser, user)
			if test.name == "User not found" {
				assert.NotNil(t, err)
			} else {
				c.Finish()
			}
		})
	}
}

