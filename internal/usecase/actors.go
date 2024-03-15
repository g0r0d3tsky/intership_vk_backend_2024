package usecase

import (
	"cinema_service/internal/domain"
	"context"
	"fmt"
	"github.com/google/uuid"
)

type ActorsRepo interface {
	CreateActor(ctx context.Context, act *domain.Actor) error
	UpdateActor(ctx context.Context, act *domain.Actor) error
	GetActors(ctx context.Context) (map[*domain.Actor][]*domain.Movie, error)
	DeleteActor(ctx context.Context, actorID uuid.UUID) error
}

type ActorsService struct {
	repo ActorsRepo
}

func NewActorsService(repo ActorsRepo) *ActorsService {
	return &ActorsService{repo: repo}
}

func (s *ActorsService) CreateActor(ctx context.Context, act *domain.Actor) error {
	err := s.repo.CreateActor(ctx, act)
	if err != nil {
		return fmt.Errorf("create actor: %w", err)
	}
	return nil
}

func (s *ActorsService) UpdateActor(ctx context.Context, act *domain.Actor) error {
	err := s.repo.UpdateActor(ctx, act)
	if err != nil {
		return fmt.Errorf("update actor: %w", err)
	}
	return nil
}

func (s *ActorsService) DeleteActor(ctx context.Context, actorID uuid.UUID) error {
	err := s.repo.DeleteActor(ctx, actorID)
	if err != nil {
		return fmt.Errorf("delete actor: %w", err)
	}
	return nil
}

func (s *ActorsService) GetActors(ctx context.Context) (map[*domain.Actor][]*domain.Movie, error) {
	actors, err := s.repo.GetActors(ctx)
	if err != nil {
		return nil, fmt.Errorf("get actors: %w", err)
	}
	return actors, nil
}
