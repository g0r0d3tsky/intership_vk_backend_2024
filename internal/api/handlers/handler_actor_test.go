package handlers 


import (
	"bytes"
	mock_service "cinema_service/internal/api/handlers/mocks"
	"cinema_service/internal/api/handlers/models"
	"cinema_service/internal/domain"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestCreateActorHandler(t *testing.T) {
	dummyError := errors.New("dummy error")
	type mockBehavior func(r *mock_service.MockActorService, actor *domain.Actor)
	testCases := []struct {
		name                 string
		inputActor           *domain.Actor
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name: "OK",
			inputActor: &domain.Actor{
				Name:      "Name",
				Surname:   "Surname",
				Sex:       "Sex",
				Birthdate: time.Time{},
			},
			mockBehavior: func(r *mock_service.MockActorService, actor *domain.Actor) {
				r.EXPECT().CreateActor(gomock.Any(), actor).Return(nil)
			},
			expectedStatusCode:   201,
			expectedResponseBody: "Actor created successfully",
		},
		{
			name: "Internal Server Error",
			inputActor: &domain.Actor{
				Name:      "Name",
				Surname:   "Surname",
				Sex:       "Sex",
				Birthdate: time.Time{},
			},
			mockBehavior: func(r *mock_service.MockActorService, actor *domain.Actor) {
				r.EXPECT().CreateActor(gomock.Any(), actor).Return(dummyError)
			},
			expectedStatusCode:   500,
			expectedResponseBody: "Failed to create actor",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			service := mock_service.NewMockActorService(c)
			tc.mockBehavior(service, tc.inputActor)

			handler := NewActorHandler(service)

			jsonData, err := json.Marshal(tc.inputActor)
			require.NoError(t, err)

			req, err := http.NewRequest("POST", "/actors", bytes.NewBuffer(jsonData))
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

			recorder := httptest.NewRecorder()

			handler.CreateActorHandler(recorder, req)

			assert.Equal(t, tc.expectedStatusCode, recorder.Code)

			if tc.expectedResponseBody != "" && tc.expectedStatusCode != 201 {
				expectedResponse := `{"error":"` + tc.expectedResponseBody + `"}`
				assert.Equal(t, expectedResponse, recorder.Body.String())
			} else if tc.expectedResponseBody != "" && tc.expectedStatusCode == 201 {
				expectedResponse := `{"status":"` + tc.expectedResponseBody + `"}`
				assert.Equal(t, expectedResponse, recorder.Body.String())
			}
		})
	}
}

func TestDeleteActorHandler(t *testing.T) {
	dummyError := errors.New("dummy error")
	type mockBehavior func(r *mock_service.MockActorService, actor *domain.Actor)
	testCases := []struct {
		name                 string
		inputActor           *domain.Actor
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name: "OK",
			inputActor: &domain.Actor{
				Name:      "Name",
				Surname:   "Surname",
				Sex:       "Sex",
				Birthdate: time.Time{},
			},
			mockBehavior: func(r *mock_service.MockActorService, actor *domain.Actor) {
				r.EXPECT().DeleteActor(gomock.Any(), actor.ID).Return(nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: "Actor deleted successfully",
		},
		{
			name: "Internal Server Error",
			inputActor: &domain.Actor{
				Name:      "Name",
				Surname:   "Surname",
				Sex:       "Sex",
				Birthdate: time.Time{},
			},
			mockBehavior: func(r *mock_service.MockActorService, actor *domain.Actor) {
				r.EXPECT().DeleteActor(gomock.Any(), actor.ID).Return(dummyError)
			},
			expectedStatusCode:   500,
			expectedResponseBody: "Failed to delete actor",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			service := mock_service.NewMockActorService(c)
			tc.mockBehavior(service, tc.inputActor)

			handler := NewActorHandler(service)

			jsonData, err := json.Marshal(tc.inputActor)
			require.NoError(t, err)

			id := "00000000-0000-0000-0000-000000000000"
			url := "/actors?id=" + url.QueryEscape(id)

			req, err := http.NewRequest("DELETE", url, bytes.NewBuffer(jsonData))
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

			recorder := httptest.NewRecorder()

			handler.DeleteActorHandler(recorder, req)

			assert.Equal(t, tc.expectedStatusCode, recorder.Code)

			if tc.expectedResponseBody != "" && tc.expectedStatusCode != 200 {
				expectedResponse := `{"error":"` + tc.expectedResponseBody + `"}`
				assert.Equal(t, expectedResponse, recorder.Body.String())
			} else if tc.expectedResponseBody != "" && tc.expectedStatusCode == 200 {
				expectedResponse := `{"status":"` + tc.expectedResponseBody + `"}`
				assert.Equal(t, expectedResponse, recorder.Body.String())
			}
		})
	}
}

func TestUpdateActorHandler(t *testing.T) {
	dummyError := errors.New("dummy error")
	type mockBehavior func(r *mock_service.MockActorService, actor *domain.Actor)
	testCases := []struct {
		name                 string
		inputActor           *domain.Actor
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name: "OK",
			inputActor: &domain.Actor{
				Name:      "Name",
				Surname:   "Surname",
				Sex:       "Sex",
				Birthdate: time.Time{},
			},
			mockBehavior: func(r *mock_service.MockActorService, actor *domain.Actor) {
				r.EXPECT().UpdateActor(gomock.Any(), actor).Return(nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: "Actor updated successfully",
		},
		{
			name: "Internal Server Error",
			inputActor: &domain.Actor{
				Name:      "Name",
				Surname:   "Surname",
				Sex:       "Sex",
				Birthdate: time.Time{},
			},
			mockBehavior: func(r *mock_service.MockActorService, actor *domain.Actor) {
				r.EXPECT().UpdateActor(gomock.Any(), actor).Return(dummyError)
			},
			expectedStatusCode:   500,
			expectedResponseBody: "Failed to update actor",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			service := mock_service.NewMockActorService(c)
			tc.mockBehavior(service, tc.inputActor)

			handler := NewActorHandler(service)

			jsonData, err := json.Marshal(tc.inputActor)
			require.NoError(t, err)

			id := "00000000-0000-0000-0000-000000000000"
			url := "/actors?id=" + url.QueryEscape(id)

			req, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonData))
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

			recorder := httptest.NewRecorder()

			handler.UpdateActorHandler(recorder, req)

			assert.Equal(t, tc.expectedStatusCode, recorder.Code)

			if tc.expectedResponseBody != "" && tc.expectedStatusCode != 200 {
				expectedResponse := `{"error":"` + tc.expectedResponseBody + `"}`
				assert.Equal(t, expectedResponse, recorder.Body.String())
			} else if tc.expectedResponseBody != "" && tc.expectedStatusCode == 200 {
				expectedResponse := `{"status":"` + tc.expectedResponseBody + `"}`
				assert.Equal(t, expectedResponse, recorder.Body.String())
			}
		})
	}
}

func TestGetActorsHandler(t *testing.T) {
	actor := &domain.Actor{
		Name:      "Name",
		Surname:   "Surname",
		Sex:       "Sex",
		Birthdate: time.Time{},
	}
	movies := []*domain.Movie{
		&domain.Movie{
			Title:       "Title",
			Description: "Description",
			Date:        time.Time{},
			Rating:      0,
		},
	}
	a := []*models.ActorMovies{
		&models.ActorMovies{
			Actor:  actor,
			Movies: movies,
		},
	}
	b := make(map[*domain.Actor][]*domain.Movie)
	b[actor] = movies
	inputActorMoviesJSON, err := json.Marshal(a)
	if err != nil {
		t.Errorf("Failed to marshal inputActorMovies to JSON: %v", err)
	}
	dummyError := errors.New("dummy error")
	type mockBehavior func(r *mock_service.MockActorService, actorMovies []*models.ActorMovies)
	testCases := []struct {
		name                 string
		inputActorMovies     []*models.ActorMovies
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody []byte
	}{
		{
			name:             "OK",
			inputActorMovies: a,
			mockBehavior: func(r *mock_service.MockActorService, actorMovies []*models.ActorMovies) {
				r.EXPECT().GetActors(gomock.Any()).Return(b, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: inputActorMoviesJSON,
		},
		{
			name:             "Internal Server Error",
			inputActorMovies: a,
			mockBehavior: func(r *mock_service.MockActorService, actorMovies []*models.ActorMovies) {
				r.EXPECT().GetActors(gomock.Any()).Return(nil, dummyError)
			},
			expectedStatusCode:   500,
			expectedResponseBody: []byte("Failed to get actors"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			service := mock_service.NewMockActorService(c)
			tc.mockBehavior(service, tc.inputActorMovies)

			handler := NewActorHandler(service)

			req, err := http.NewRequest("GET", "/actors", nil)
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

			recorder := httptest.NewRecorder()

			handler.GetActorsHandler(recorder, req)

			assert.Equal(t, tc.expectedStatusCode, recorder.Code)

			if tc.expectedStatusCode != 200 {
				assert.Equal(t, `{"error":"`+string(tc.expectedResponseBody)+`"}`, recorder.Body.String())
			} else if tc.expectedStatusCode == 200 {
				assert.JSONEq(t, string(tc.expectedResponseBody), recorder.Body.String(), "JSON strings do not match")
			}
		})
	}
}