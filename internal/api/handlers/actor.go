package handlers

import (
	"cinema_service/internal/api/handlers/models"
	"cinema_service/internal/domain"
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"net/http"
)

type ActorService interface {
	CreateActor(ctx context.Context, act *domain.Actor) error
	UpdateActor(ctx context.Context, act *domain.Actor) error
	DeleteActor(ctx context.Context, actorID uuid.UUID) error
	GetActors(ctx context.Context) (map[*domain.Actor][]*domain.Movie, error)
}

type ActorHandler struct {
	service ActorService
}

func NewActorHandler(service ActorService) *ActorHandler {
	return &ActorHandler{service: service}
}

func (h *ActorHandler) CreateActorHandler(w http.ResponseWriter, r *http.Request) {
	var input models.Actor
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	actor := &domain.Actor{
		Name:    input.Name,
		Surname: input.Surname,
		Sex:     input.Sex,
	}
	err = h.service.CreateActor(r.Context(), actor)
	if err != nil {
		http.Error(w, "Failed to create actor", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *ActorHandler) UpdateActorHandler(w http.ResponseWriter, r *http.Request) {
	actorIDStr := r.URL.Query().Get("actor_id")
	_, err := uuid.Parse(actorIDStr)
	if err != nil {
		http.Error(w, "Invalid actor ID", http.StatusBadRequest)
		return
	}

	var input models.Actor
	err = json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	actor := &domain.Actor{
		Name:    input.Name,
		Surname: input.Surname,
		Sex:     input.Sex,
	}
	err = h.service.UpdateActor(r.Context(), actor)
	if err != nil {
		http.Error(w, "Failed to update actor", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *ActorHandler) DeleteActorHandler(w http.ResponseWriter, r *http.Request) {
	actorIDStr := r.URL.Query().Get("actor_id")
	actorID, err := uuid.Parse(actorIDStr)
	if err != nil {
		http.Error(w, "Invalid actor ID", http.StatusBadRequest)
		return
	}

	err = h.service.DeleteActor(r.Context(), actorID)
	if err != nil {
		http.Error(w, "Failed to delete actor", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *ActorHandler) GetActorsHandler(w http.ResponseWriter, r *http.Request) {
	actors, err := h.service.GetActors(r.Context())
	if err != nil {
		http.Error(w, "Failed to get actors", http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(actors)
	if err != nil {
		return
	}
}
