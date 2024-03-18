package handlers

import (
	"cinema_service/internal/api/handlers/models"
	"cinema_service/internal/domain"
	"context"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
)

//go:generate mockgen -source=actor.go -destination=mocks/actorServiceMock.go

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
// @Security ApiKeyAuth
// @Param actor body models.Actor true "Actor object"
// @Success 201 {object} statusResponse "Actor created successfully"
// @Failure 400 {object} errorResponse "Invalid request payload"
// @Failure 500 {object} errorResponse "Failed to create actor"
// @Router /actors [post]
func (h *ActorHandler) CreateActorHandler(w http.ResponseWriter, r *http.Request) {
	var input models.Actor
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		NewErrorResponse(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	actor := &domain.Actor{
		Name:    input.Name,
		Surname: input.Surname,
		Sex:     input.Sex,
	}
	err = h.service.CreateActor(r.Context(), actor)
	if err != nil {
		NewErrorResponse(w, http.StatusInternalServerError, "Failed to create actor")
		return
	}

	sendJSONResponse(w, http.StatusCreated, statusResponse{
		Status: "Actor created successfully",
	})
}

// UpdateActorHandler updates actor information.
// @Summary Update actor information
// @Description Updates actor information based on the input data.
// @Tags Actors
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id query string true "Actor ID"
// @Param actor body models.Actor true "Updated actor information"
// @Success 200 {object} statusResponse
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /actors [put]
func (h *ActorHandler) UpdateActorHandler(w http.ResponseWriter, r *http.Request) {
	actorIDStr := r.URL.Query().Get("id")
	if actorIDStr == "" {
		NewErrorResponse(w, http.StatusBadRequest, "Actor ID parameter is required")
		return
	}

	id, err := uuid.Parse(actorIDStr)
	if err != nil {
		NewErrorResponse(w, http.StatusBadRequest,
			"Invalid actor ID",
		)
		return
	}

	var input models.Actor
	err = json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		NewErrorResponse(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	actor := &domain.Actor{
		ID:      id,
		Name:    input.Name,
		Surname: input.Surname,
		Sex:     input.Sex,
	}

	err = h.service.UpdateActor(r.Context(), actor)
	if err != nil {
		NewErrorResponse(w, http.StatusInternalServerError, "Failed to update actor")
		return
	}

	sendJSONResponse(w, http.StatusOK, statusResponse{
		Status: "Actor updated successfully",
	})
}

// DeleteActorHandler deletes an actor.
// @Summary Delete an actor
// @Description Deletes an actor based on the provided actor ID.
// @Tags Actors
// @Security ApiKeyAuth
// @Param id query string true "Actor ID"
// @Success 200 {object} statusResponse
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /actors [delete]
func (h *ActorHandler) DeleteActorHandler(w http.ResponseWriter, r *http.Request) {
	actorIDStr := r.URL.Query().Get("id")
	if actorIDStr == "" {
		NewErrorResponse(w, http.StatusBadRequest, "Actor ID parameter is required")
		return
	}

	actorID, err := uuid.Parse(actorIDStr)
	if err != nil {
		NewErrorResponse(w, http.StatusBadRequest, "Invalid actor ID")
		return
	}

	err = h.service.DeleteActor(r.Context(), actorID)
	if err != nil {
		NewErrorResponse(w, http.StatusInternalServerError, "Failed to delete actor")
		return
	}

	sendJSONResponse(w, http.StatusOK, statusResponse{
		Status: "Actor deleted successfully",
	})
}

// GetActorsHandler retrieves a list of actors.
// @Summary Get Actors
// @Description Retrieves a list of actors
// @Tags Actors
// @Accept json
// @Param id path int true "Actor ID"
// @Security ApiKeyAuth
// @Success 200 {object} models.ActorMovies "Actors retrieved successfully"
// @Failure 500 {object} errorResponse
// @Router /actors [get]
func (h *ActorHandler) GetActorsHandler(w http.ResponseWriter, r *http.Request) {
	actors, err := h.service.GetActors(r.Context())
	if err != nil {
		NewErrorResponse(w, http.StatusInternalServerError, "Failed to get actors")
		return
	}

	actorMoviesList := make([]*models.ActorMovies, 0, len(actors))

	for actor, movies := range actors {
		actorMovies := &models.ActorMovies{
			Actor:  actor,
			Movies: movies,
		}
		actorMoviesList = append(actorMoviesList, actorMovies)
	}

	sendJSONResponse(w, http.StatusOK, actorMoviesList)
}

// TODO: authorization
func (h *ActorHandler) RegisterActor(mux *http.ServeMux,
	authentication Middleware, authorization Middleware) *http.ServeMux {
	mux.HandleFunc("GET /api/v1/actors", authentication(h.GetActorsHandler))
	mux.HandleFunc("POST /api/v1/actors", authentication(authorization(h.CreateActorHandler)))
	mux.HandleFunc("PUT /api/v1/actors", authentication(authorization(h.UpdateActorHandler)))
	mux.HandleFunc("DELETE /api/v1/actors", authentication(authorization(h.DeleteActorHandler)))
	return mux
}

//TODO: пофикси удаление из таблиц
