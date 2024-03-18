package usecase

import (
	"context"
	"fmt"
	"testing"

	"cinema_service/internal/domain"
	mock_repo "cinema_service/internal/usecase/mocks"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)


func TestCreateActor(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repo.NewMockActorsRepo(ctrl) 
	actorService := NewActorsService(mockRepo)

	testCases := []struct {
		name     string
		actor    *domain.Actor
		mockFunc func()
		wantErr  bool
	}{
		{
			name:  "Create actor successfully",
			actor: &domain.Actor{},
			mockFunc: func() {
				mockRepo.EXPECT().CreateActor(gomock.Any(), gomock.Any()).Return(nil)
			},
			wantErr: false,
		},
		{
			name:  "Create actor fails",
			actor: &domain.Actor{},
			mockFunc: func() {
				mockRepo.EXPECT().CreateActor(gomock.Any(), gomock.Any()).Return(fmt.Errorf("some error"))
			},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockFunc() 

			err := actorService.CreateActor(context.Background(), tc.actor)

			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
func TestUpdateActor(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repo.NewMockActorsRepo(ctrl)

	actorService := NewActorsService(mockRepo)

	testCases := []struct {
		name     string
		actor    *domain.Actor
		mockFunc func()
		wantErr  bool
	}{
		{
			name:  "Update actor successfully",
			actor: &domain.Actor{},
			mockFunc: func() {
				mockRepo.EXPECT().UpdateActor(gomock.Any(), gomock.Any()).Return(nil)
			},
			wantErr: false,
		},
		{
			name:  "Update actor fails",
			actor: &domain.Actor{},
			mockFunc: func() {
				mockRepo.EXPECT().UpdateActor(gomock.Any(), gomock.Any()).Return(fmt.Errorf("some error"))
			},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockFunc()

			err := actorService.UpdateActor(context.Background(), tc.actor)

			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestDeleteActor(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repo.NewMockActorsRepo(ctrl)

	actorService := NewActorsService(mockRepo)

	testCases := []struct {
		name     string
		actorID  uuid.UUID
		mockFunc func()
		wantErr  bool
	}{
		{
			name:    "Delete actor successfully",
			actorID: uuid.UUID{},
			mockFunc: func() {
				mockRepo.EXPECT().DeleteActor(gomock.Any(), gomock.Any()).Return(nil)
			},
			wantErr: false,
		},
		{
			name:    "Delete actor fails",
			actorID: uuid.UUID{},
			mockFunc: func() {
				mockRepo.EXPECT().DeleteActor(gomock.Any(), gomock.Any()).Return(fmt.Errorf("some error"))
			},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockFunc()

			err := actorService.DeleteActor(context.Background(), tc.actorID)

			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestGetActors(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repo.NewMockActorsRepo(ctrl)

	actorService := NewActorsService(mockRepo)

	testCases := []struct {
		name     string
		mockFunc func()
		wantErr  bool
	}{
		{
			name: "Get actors successfully",
			mockFunc: func() {
				mockRepo.EXPECT().GetActors(gomock.Any()).Return(make(map[*domain.Actor][]*domain.Movie), nil)
			},
			wantErr: false,
		},
		{
			name: "Get actors fails",
			mockFunc: func() {
				mockRepo.EXPECT().GetActors(gomock.Any()).Return(nil, fmt.Errorf("some error"))
			},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockFunc()

			_, err := actorService.GetActors(context.Background())

			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}