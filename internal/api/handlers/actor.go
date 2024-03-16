package handlers

import (
	"cinema_service/internal/api/handlers/models"
	"cinema_service/internal/domain"
	"context"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
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

type Middleware func(handlerFunc http.HandlerFunc) http.HandlerFunc

func NewActorHandler(service ActorService) *ActorHandler {
	return &ActorHandler{service: service}
}

// CreateActorHandler creates a new actor.
// @Summary Create Actor
// @Description Creates a new actor
// @Tags Actors
// @Accept json
// @Produce json
// @Param body models.Actor true "Actor object"
// @Success 201 {string} string "Actor created successfully"
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /actors [post]
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

// UpdateActorHandler updates an existing actor.
// @Summary Update Actor
// @Description Updates an existing actor
// @Tags Actors
// @Accept json
// @Produce json
// @Param actor_id query string true "Actor ID"
// @Param body  models.Actor true "Actor object"
// @Success 200 {string} string "Actor updated successfully"
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /actors [put]
func (h *ActorHandler) UpdateActorHandler(w http.ResponseWriter, r *http.Request) {
	actorIDStr := r.URL.Query().Get("id")
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

// DeleteActorHandler deletes an actor.
// @Summary Delete Actor
// @Description Deletes an actor
// @Tags Actors
// @Accept json
// @Produce json
// @Param actor_id query string true "Actor ID"
// @Success 200 {string} string "Actor deleted successfully"
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /actors [delete]
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

// GetActorsHandler retrieves a list of actors.
// @Summary Get Actors
// @Description Retrieves a list of actors
// @Tags Actors
// @Accept json
// @Produce json
// @Success 200 {object} map[*domain.Actor][]*domain.Movie "Actors retrieved successfully"
// @Failure 500 {object} ErrorResponse
// @Router /actors [get]
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

// TODO: authorization
func (h *ActorHandler) RegisterActor(mux *http.ServeMux,
	authentication Middleware, authorization Middleware) *http.ServeMux {
	mux.HandleFunc("GET /api/v1/actors/", authentication(h.GetActorsHandler))
	mux.HandleFunc("POST /api/v1/actors/", authentication(authorization(h.CreateActorHandler)))
	mux.HandleFunc("PUT /api/v1/actors/{id}", authentication(authorization(h.UpdateActorHandler)))
	mux.HandleFunc("DELETE /api/v1/actors/{id}", authentication(authentication(h.DeleteActorHandler)))
	return mux
}
