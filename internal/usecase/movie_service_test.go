package usecase

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"cinema_service/internal/domain"
	mock_repo "cinema_service/internal/usecase/mocks"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestCreateMovie(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repo.NewMockMovieRepo(ctrl)

	movieService := NewMovieService(mockRepo)

	testCases := []struct {
		name     string
		movie    *domain.Movie
		mockFunc func()
		wantErr  bool
	}{
		{
			name:  "Create movie successfully",
			movie: &domain.Movie{},
			mockFunc: func() {
				mockRepo.EXPECT().CreateMovie(gomock.Any(), gomock.Any()).Return(nil)
			},
			wantErr: false,
		},
		{
			name:  "Create movie fails",
			movie: &domain.Movie{},
			mockFunc: func() {
				mockRepo.EXPECT().CreateMovie(gomock.Any(), gomock.Any()).Return(fmt.Errorf("some error"))
			},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockFunc()

			err := movieService.CreateMovie(context.Background(), tc.movie)

			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestUpdateMovie(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repo.NewMockMovieRepo(ctrl)

	movieService := NewMovieService(mockRepo)

	testCases := []struct {
		name     string
		movie    *domain.Movie
		mockFunc func()
		wantErr  bool
	}{
		{
			name:  "Update movie successfully",
			movie: &domain.Movie{},
			mockFunc: func() {
				mockRepo.EXPECT().UpdateMovie(gomock.Any(), gomock.Any()).Return(nil)
			},
			wantErr: false,
		},
		{
			name:  "Update movie fails",
			movie: &domain.Movie{},
			mockFunc: func() {
				mockRepo.EXPECT().UpdateMovie(gomock.Any(), gomock.Any()).Return(fmt.Errorf("some error"))
			},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockFunc()

			err := movieService.UpdateMovie(context.Background(), tc.movie)

			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestDeleteMovie(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repo.NewMockMovieRepo(ctrl)

	movieService := NewMovieService(mockRepo)

	testCases := []struct {
		name     string
		movieID  uuid.UUID
		mockFunc func()
		wantErr  bool
	}{
		{
			name:    "Delete movie successfully",
			movieID: uuid.New(),
			mockFunc: func() {
				mockRepo.EXPECT().DeleteMovie(gomock.Any(), gomock.Any()).Return(nil)
			},
			wantErr: false,
		},
		{
			name:    "Delete movie fails",
			movieID: uuid.New(),
			mockFunc: func() {
				mockRepo.EXPECT().DeleteMovie(gomock.Any(), gomock.Any()).Return(fmt.Errorf("some error"))
			},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockFunc()

			err := movieService.DeleteMovie(context.Background(), tc.movieID)

			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestGetMovies(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repo.NewMockMovieRepo(ctrl)

	movieService := NewMovieService(mockRepo)

	testCases := []struct {
		name     string
		mockFunc func()
		want     []*domain.Movie
		wantErr  bool
	}{
		{
			name: "Get movies successfully",
			mockFunc: func() {
				mockRepo.EXPECT().GetMovies(gomock.Any()).Return([]*domain.Movie{
					&domain.Movie{ID: func() uuid.UUID { id, _ := uuid.Parse("6ec91a6d-12ce-4bd1-b7f1-e70b94eeef0b"); return id }(), Title: "Movie 1"},
					&domain.Movie{ID: func() uuid.UUID { id, _ := uuid.Parse("ae80bb3c-a74c-4335-9560-6e0e4df63c29"); return id }(), Title: "Movie 2"},
				}, nil)
			},
			want: []*domain.Movie{
				&domain.Movie{ID: func() uuid.UUID { id, _ := uuid.Parse("6ec91a6d-12ce-4bd1-b7f1-e70b94eeef0b"); return id }(), Title: "Movie 1"},
				&domain.Movie{ID: func() uuid.UUID { id, _ := uuid.Parse("ae80bb3c-a74c-4335-9560-6e0e4df63c29"); return id }(), Title: "Movie 2"},
			},
			wantErr: false,
		},
		{
			name: "Get movies fails",
			mockFunc: func() {
				mockRepo.EXPECT().GetMovies(gomock.Any()).Return(nil, fmt.Errorf("some error"))
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockFunc()

			movies, err := movieService.GetMovies(context.Background())

			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tc.want, movies)
		})
	}
}
func TestGetMoviesFilter(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repo.NewMockMovieRepo(ctrl)

	movieService := NewMovieService(mockRepo)

	movies := []*domain.Movie{
		&domain.Movie{ID: func() uuid.UUID { id, _ := uuid.Parse("6ec91a6d-12ce-4bd1-b7f1-e70b94eeef0b"); return id }(), Title: "Movie B", Rating: 8.5, Date: time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)},
		&domain.Movie{ID: func() uuid.UUID { id, _ := uuid.Parse("ae80bb3c-a74c-4335-9560-6e0e4df63c29"); return id }(), Title: "Movie A", Rating: 7.5, Date: time.Date(2021, 12, 31, 0, 0, 0, 0, time.UTC)},
		&domain.Movie{ID: func() uuid.UUID { id, _ := uuid.Parse("c997cb73-b393-4634-aed0-5c3c2137324d"); return id }(), Title: "Movie C", Rating: 9.5, Date: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)},
		&domain.Movie{ID: func() uuid.UUID { id, _ := uuid.Parse("12a0024e-1219-4a4a-817b-2b9e8ff7616f"); return id }(), Title: "Movie D", Rating: 4.5, Date: time.Date(1971, 12, 31, 0, 0, 0, 0, time.UTC)},
	}
	testCases := []struct {
		name     string
		filter   string
		movies   []*domain.Movie
		mockFunc func()
		want     []*domain.Movie
		wantErr  bool
	}{
		{
			name:   "Filter by title correct",
			filter: "title",
			movies: movies,
			mockFunc: func() {
				mockRepo.EXPECT().GetMovies(gomock.Any()).Return(movies, nil)
			},
			want: []*domain.Movie{
				&domain.Movie{ID: func() uuid.UUID { id, _ := uuid.Parse("ae80bb3c-a74c-4335-9560-6e0e4df63c29"); return id }(), Title: "Movie A", Rating: 7.5, Date: time.Date(2021, 12, 31, 0, 0, 0, 0, time.UTC)},
				&domain.Movie{ID: func() uuid.UUID { id, _ := uuid.Parse("6ec91a6d-12ce-4bd1-b7f1-e70b94eeef0b"); return id }(), Title: "Movie B", Rating: 8.5, Date: time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)},
				&domain.Movie{ID: func() uuid.UUID { id, _ := uuid.Parse("c997cb73-b393-4634-aed0-5c3c2137324d"); return id }(), Title: "Movie C", Rating: 9.5, Date: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)},
				&domain.Movie{ID: func() uuid.UUID { id, _ := uuid.Parse("12a0024e-1219-4a4a-817b-2b9e8ff7616f"); return id }(), Title: "Movie D", Rating: 4.5, Date: time.Date(1971, 12, 31, 0, 0, 0, 0, time.UTC)},
			},
			wantErr: false,
		},
		{
			name:   "Filter by rating correct",
			filter: "rating",
			movies: movies,
			mockFunc: func() {
				mockRepo.EXPECT().GetMovies(gomock.Any()).Return(movies, nil)
			},
			want: []*domain.Movie{
					&domain.Movie{ID: func() uuid.UUID { id, _ := uuid.Parse("c997cb73-b393-4634-aed0-5c3c2137324d"); return id }(), Title: "Movie C", Rating: 9.5, Date: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)},
					&domain.Movie{ID: func() uuid.UUID { id, _ := uuid.Parse("6ec91a6d-12ce-4bd1-b7f1-e70b94eeef0b"); return id }(), Title: "Movie B", Rating: 8.5, Date: time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)},
					&domain.Movie{ID: func() uuid.UUID { id, _ := uuid.Parse("ae80bb3c-a74c-4335-9560-6e0e4df63c29"); return id }(), Title: "Movie A", Rating: 7.5, Date: time.Date(2021, 12, 31, 0, 0, 0, 0, time.UTC)},
					&domain.Movie{ID: func() uuid.UUID { id, _ := uuid.Parse("12a0024e-1219-4a4a-817b-2b9e8ff7616f"); return id }(), Title: "Movie D", Rating: 4.5, Date: time.Date(1971, 12, 31, 0, 0, 0, 0, time.UTC)},
			},
			wantErr: false,
	
		},
		{
			name:   "Filter by date correct",
			filter: "date",
			movies: movies,
			mockFunc: func() {
				mockRepo.EXPECT().GetMovies(gomock.Any()).Return(movies, nil)
			},
			want: []*domain.Movie{
					&domain.Movie{ID: func() uuid.UUID { id, _ := uuid.Parse("c997cb73-b393-4634-aed0-5c3c2137324d"); return id }(), Title: "Movie C", Rating: 9.5, Date: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)},
					&domain.Movie{ID: func() uuid.UUID { id, _ := uuid.Parse("6ec91a6d-12ce-4bd1-b7f1-e70b94eeef0b"); return id }(), Title: "Movie B", Rating: 8.5, Date: time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)},
					&domain.Movie{ID: func() uuid.UUID { id, _ := uuid.Parse("ae80bb3c-a74c-4335-9560-6e0e4df63c29"); return id }(), Title: "Movie A", Rating: 7.5, Date: time.Date(2021, 12, 31, 0, 0, 0, 0, time.UTC)},
					&domain.Movie{ID: func() uuid.UUID { id, _ := uuid.Parse("12a0024e-1219-4a4a-817b-2b9e8ff7616f"); return id }(), Title: "Movie D", Rating: 4.5, Date: time.Date(1971, 12, 31, 0, 0, 0, 0, time.UTC)},
			},
			wantErr: false,
	
		},
		{
			name:   "Filter empty",
			filter: "",
			movies: movies,
			mockFunc: func() {
				mockRepo.EXPECT().GetMovies(gomock.Any()).Return(movies, nil)
			},
			want: []*domain.Movie{
				&domain.Movie{ID: func() uuid.UUID { id, _ := uuid.Parse("c997cb73-b393-4634-aed0-5c3c2137324d"); return id }(), Title: "Movie C", Rating: 9.5, Date: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)},
				&domain.Movie{ID: func() uuid.UUID { id, _ := uuid.Parse("6ec91a6d-12ce-4bd1-b7f1-e70b94eeef0b"); return id }(), Title: "Movie B", Rating: 8.5, Date: time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)},
				&domain.Movie{ID: func() uuid.UUID { id, _ := uuid.Parse("ae80bb3c-a74c-4335-9560-6e0e4df63c29"); return id }(), Title: "Movie A", Rating: 7.5, Date: time.Date(2021, 12, 31, 0, 0, 0, 0, time.UTC)},
				&domain.Movie{ID: func() uuid.UUID { id, _ := uuid.Parse("12a0024e-1219-4a4a-817b-2b9e8ff7616f"); return id }(), Title: "Movie D", Rating: 4.5, Date: time.Date(1971, 12, 31, 0, 0, 0, 0, time.UTC)},
			},
			wantErr: false,
	
		},
		{
			name:   "Filter invalid",
			filter: "invalid",
			movies: movies,
			mockFunc: func() {
				mockRepo.EXPECT().GetMovies(gomock.Any()).Return(movies, nil)
			},
			want: nil,
			wantErr: true,
	
		},
		{
			name:    "Repo_error",
			filter:  "invalid",
			mockFunc: func() {
				mockRepo.EXPECT().GetMovies(gomock.Any()).Return(nil, errors.New("get movies: dummyError"))
			},
			want:    nil,
			wantErr: true,
	
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockFunc()

			movies, err := movieService.GetMoviesFilter(context.Background(), tc.filter)

			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
	
			assert.Equal(t, tc.want, movies)
		})
	}

}
func TestGetMoviesBySnippet(t *testing.T) {
	type mockBehavior func(r *mock_repo.MockMovieRepo, snippet string)

	testCases := []struct {
		name            string
		snippet         string
		mockBehavior    mockBehavior
		expectedMovies  []*domain.Movie
		expectedErr     string
	}{
		{
			name:    "Success",
			snippet: "action",
			mockBehavior: func(r *mock_repo.MockMovieRepo, snippet string) {
				r.EXPECT().GetMoviesBySnippet(gomock.Any(), snippet).Return([]*domain.Movie{
					{Title: "Movie 1"},
					{Title: "Movie 2"},
				}, nil)
			},
			expectedMovies: []*domain.Movie{
				{Title: "Movie 1"},
				{Title: "Movie 2"},
			},
			expectedErr: "",
		},
		{
			name:    "Repository error",
			snippet: "drama",
			mockBehavior: func(r *mock_repo.MockMovieRepo, snippet string) {
				r.EXPECT().GetMoviesBySnippet(gomock.Any(), snippet).Return(nil, errors.New("repository error"))
			},
			expectedMovies: nil,
			expectedErr:    "get movies by snippet: repository error",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := mock_repo.NewMockMovieRepo(ctrl)
			tc.mockBehavior(mockRepo, tc.snippet)

			service := NewMovieService(mockRepo)

			movies, err := service.GetMoviesBySnippet(context.Background(), tc.snippet)

			if tc.expectedErr == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.expectedErr)
			}
			assert.Equal(t, tc.expectedMovies, movies)
		})
	}
}