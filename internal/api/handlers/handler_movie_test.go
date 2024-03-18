package handlers

import (
	"bytes"
	mock_service "cinema_service/internal/api/handlers/mocks"
	"cinema_service/internal/domain"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

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

