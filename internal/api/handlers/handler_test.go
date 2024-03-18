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

func TestCreateMovieHandler(t *testing.T) {
	dummyError := errors.New("dummy error")
	type mockBehavior func(r *mock_service.MockMovieService, movie *domain.Movie)
	testCases := []struct {
		name string

		inputMovie           *domain.Movie
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name: "OK",
			inputMovie: &domain.Movie{
				Title:       "Test Movie",
				Description: "Test Description",
				Rating:      4.5,
			},
			mockBehavior: func(r *mock_service.MockMovieService, movie *domain.Movie) {
				r.EXPECT().CreateMovie(gomock.Any(), movie).Return(nil)
			},
			expectedStatusCode:   201,
			expectedResponseBody: "",
		},
		{
			name: "Internal Server Error",
			inputMovie: &domain.Movie{
				Title:       "Test Movie",
				Description: "Test Description",
				Rating:      4.5,
			},
			mockBehavior: func(r *mock_service.MockMovieService, movie *domain.Movie) {
				r.EXPECT().CreateMovie(gomock.Any(), movie).Return(errors.New(dummyError.Error()))
			},
			expectedStatusCode:   500,
			expectedResponseBody: "Failed to create movie",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			service := mock_service.NewMockMovieService(c)
			tc.mockBehavior(service, tc.inputMovie)

			handler := NewMovieHandler(service)

			jsonData, err := json.Marshal(tc.inputMovie)
			require.NoError(t, err)

			req, err := http.NewRequest("POST", "/movies", bytes.NewBuffer(jsonData))
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

			recorder := httptest.NewRecorder()

			handler.CreateMovieHandler(recorder, req)

			assert.Equal(t, tc.expectedStatusCode, recorder.Code)

			if tc.expectedResponseBody != "" {
				expectedResponse := `{"error":"` + tc.expectedResponseBody + `"}`
				assert.Equal(t, expectedResponse, recorder.Body.String())
			}
		})
	}
}
func TestUpdateMovieHandler(t *testing.T) {
	dummyError := errors.New("dummy error")
	type mockBehavior func(r *mock_service.MockMovieService, movie *domain.Movie)
	testCases := []struct {
		name                 string
		inputMovie           *domain.Movie
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name: "OK",
			inputMovie: &domain.Movie{
				Title:       "Test Movie",
				Description: "Test Description",
				Rating:      4.5,
			},
			mockBehavior: func(r *mock_service.MockMovieService, movie *domain.Movie) {
				r.EXPECT().UpdateMovie(gomock.Any(), movie).Return(nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: "",
		},
		{
			name: "Internal Server Error",
			inputMovie: &domain.Movie{
				Title:       "Test Movie",
				Description: "Test Description",
				Rating:      4.5,
			},
			mockBehavior: func(r *mock_service.MockMovieService, movie *domain.Movie) {
				r.EXPECT().UpdateMovie(gomock.Any(), movie).Return(errors.New(dummyError.Error()))
			},
			expectedStatusCode:   500,
			expectedResponseBody: "Failed to update movie",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			service := mock_service.NewMockMovieService(c)
			tc.mockBehavior(service, tc.inputMovie)

			handler := NewMovieHandler(service)

			jsonData, err := json.Marshal(tc.inputMovie)
			require.NoError(t, err)

			id := "00000000-0000-0000-0000-000000000000"
			url := "/movies?id=" + url.QueryEscape(id)
			req, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonData))
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

			recorder := httptest.NewRecorder()

			handler.UpdateMovieHandler(recorder, req)

			assert.Equal(t, tc.expectedStatusCode, recorder.Code)

			if tc.expectedResponseBody != "" && tc.expectedStatusCode != 200 {
				expectedResponse := `{"error":"` + tc.expectedResponseBody + `"}`
				assert.Equal(t, expectedResponse, recorder.Body.String())
			} else if tc.expectedResponseBody != "" && tc.expectedStatusCode == 200 {
				expectedResponse := `{"Status":"` + tc.expectedResponseBody + `"}`
				assert.Equal(t, expectedResponse, recorder.Body.String())
			}
		})
	}
}

func TestDeleteMovieHandler(t *testing.T) {
	dummyError := errors.New("dummy error")
	type mockBehavior func(r *mock_service.MockMovieService, movie *domain.Movie)
	testCases := []struct {
		name                 string
		inputMovie           *domain.Movie
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name: "OK",
			inputMovie: &domain.Movie{
				Title:       "Test Movie",
				Description: "Test Description",
				Rating:      4.5,
			},
			mockBehavior: func(r *mock_service.MockMovieService, movie *domain.Movie) {
				r.EXPECT().DeleteMovie(gomock.Any(), movie.ID).Return(nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: "Movie deleted successfully",
		},
		{
			name: "Internal Server Error",
			inputMovie: &domain.Movie{
				Title:       "Test Movie",
				Description: "Test Description",
				Rating:      4.5,
			},
			mockBehavior: func(r *mock_service.MockMovieService, movie *domain.Movie) {
				r.EXPECT().DeleteMovie(gomock.Any(), movie.ID).Return(errors.New(dummyError.Error()))
			},
			expectedStatusCode:   500,
			expectedResponseBody: "Failed to delete movie",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			service := mock_service.NewMockMovieService(c)
			tc.mockBehavior(service, tc.inputMovie)

			handler := NewMovieHandler(service)

			jsonData, err := json.Marshal(tc.inputMovie)
			require.NoError(t, err)

			id := "00000000-0000-0000-0000-000000000000"
			url := "/movies?id=" + url.QueryEscape(id)
			req, err := http.NewRequest("DELETE", url, bytes.NewBuffer(jsonData))
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

			recorder := httptest.NewRecorder()

			handler.DeleteMovieHandler(recorder, req)

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

func TestGetMoviesFilterHandler(t *testing.T) {
	dummyError := errors.New("dummy error")
	type mockBehavior func(r *mock_service.MockMovieService, movies []*domain.Movie)
	testCases := []struct {
		name                 string
		inputMovies          []*domain.Movie
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name: "OK",
			inputMovies: []*domain.Movie{
				&domain.Movie{
					Title:       "Test Movie",
					Description: "Test Description",
					Rating:      4.5,
				},
			},
			mockBehavior: func(r *mock_service.MockMovieService, movies []*domain.Movie) {
				r.EXPECT().GetMoviesFilter(gomock.Any(), "rating").Return(movies, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: "",
		},
		{
			name: "OK",
			inputMovies: []*domain.Movie{
				&domain.Movie{
					Title:       "Test Movie",
					Description: "Test Description",
					Rating:      4.5,
				},
			},
			mockBehavior: func(r *mock_service.MockMovieService, movies []*domain.Movie) {
				r.EXPECT().GetMoviesFilter(gomock.Any(), "rating").Return(nil, dummyError)
			},
			expectedStatusCode:   500,
			expectedResponseBody: "Failed to get movies",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			service := mock_service.NewMockMovieService(c)
			tc.mockBehavior(service, tc.inputMovies)

			handler := NewMovieHandler(service)

			jsonData, err := json.Marshal(tc.inputMovies)
			require.NoError(t, err)

			rating := "rating"
			url := "/movies?filter=" + url.QueryEscape(rating)
			req, err := http.NewRequest("GET", url, bytes.NewBuffer(jsonData))
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

			recorder := httptest.NewRecorder()

			handler.GetMoviesFilterHandler(recorder, req)

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

func TestGetMoviesBySnippetHandler(t *testing.T) {
	dummyError := errors.New("dummy error")
	type mockBehavior func(r *mock_service.MockMovieService, movies []*domain.Movie)
	testCases := []struct {
		name                 string
		inputMovies          []*domain.Movie
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name: "OK",
			inputMovies: []*domain.Movie{
				&domain.Movie{
					Title:       "Test Movie",
					Description: "Test Description",
					Rating:      4.5,
				},
			},
			mockBehavior: func(r *mock_service.MockMovieService, movies []*domain.Movie) {
				r.EXPECT().GetMoviesBySnippet(gomock.Any(), "VK is one of the best places to start career").Return(movies, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: "",
		},
		{
			name: "OK",
			inputMovies: []*domain.Movie{
				&domain.Movie{
					Title:       "Test Movie",
					Description: "Test Description",
					Rating:      4.5,
				},
			},
			mockBehavior: func(r *mock_service.MockMovieService, movies []*domain.Movie) {
				r.EXPECT().GetMoviesBySnippet(gomock.Any(), "VK is one of the best places to start career").Return(nil, dummyError)
			},
			expectedStatusCode:   500,
			expectedResponseBody: "Failed to get movies",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			service := mock_service.NewMockMovieService(c)
			tc.mockBehavior(service, tc.inputMovies)

			handler := NewMovieHandler(service)

			jsonData, err := json.Marshal(tc.inputMovies)
			require.NoError(t, err)

			rating := "VK is one of the best places to start career"
			url := "/movies?snippet=" + url.QueryEscape(rating)
			req, err := http.NewRequest("GET", url, bytes.NewBuffer(jsonData))
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

			recorder := httptest.NewRecorder()

			handler.GetMoviesBySnippetHandler(recorder, req)

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

func TestSignInHandler(t *testing.T) {
	dummyError := errors.New("dummy error")
	type mockBehavior func(r *mock_service.MockUserService, input signInInput)
	testCases := []struct {
		name                 string
		input                signInInput
		token                string
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name: "OK",
			input: signInInput{
				Login:    "login",
				Password: "password",
			},
			mockBehavior: func(r *mock_service.MockUserService, input signInInput) {
				r.EXPECT().GenerateToken(gomock.Any(), input.Login, input.Password).Return("token", nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: "token",
		},
		{
			name: "Internal Server Error",
			input: signInInput{
				Login:    "login",
				Password: "password",
			},
			mockBehavior: func(r *mock_service.MockUserService, input signInInput) {
				r.EXPECT().GenerateToken(gomock.Any(), input.Login, input.Password).Return("", dummyError)
			},
			expectedStatusCode:   500,
			expectedResponseBody: "Generating Token error",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			service := mock_service.NewMockUserService(c)
			tc.mockBehavior(service, tc.input)

			handler := NewUserHandler(service)

			jsonData, err := json.Marshal(tc.input)
			require.NoError(t, err)

			req, err := http.NewRequest("PUT", "/signIn", bytes.NewBuffer(jsonData))
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

			recorder := httptest.NewRecorder()
			handler.SignIn(recorder, req)

			assert.Equal(t, tc.expectedStatusCode, recorder.Code)

			if tc.expectedStatusCode != 200 {
				expectedResponse := `{"error":"` + tc.expectedResponseBody + `"}`
				assert.Equal(t, expectedResponse, recorder.Body.String())
			}
			if tc.expectedStatusCode == 200 {
				expectedResponse := `{"token":"` + tc.expectedResponseBody + `"}`
				assert.JSONEq(t, expectedResponse, recorder.Body.String())
			}
		})
	}
}
