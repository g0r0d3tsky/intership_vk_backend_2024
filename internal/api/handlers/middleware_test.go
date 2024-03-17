package handlers

import (
	mock_service "cinema_service/internal/api/handlers/mocks"
	"cinema_service/internal/api/middleware"
	"cinema_service/internal/usecase"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"

	"go.uber.org/mock/gomock"
)

func TestHandlerUserIdentity(t *testing.T) {
	type mockBehavior func(s *mock_service.MockUserService, token string)

	tests := []struct {
		name                 string
		headerName           string
		headerValue          string
		token                string
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:        "OK",
			headerName:  "Authorization",
			headerValue: "Bearer valid_token",
			token:       "valid_token",
			mockBehavior: func(s *mock_service.MockUserService, token string) {
				s.EXPECT().ParseToken(token).Return(&domain.User{UserID: "00000000-0000-0000-0000-000000000000", Role: "user"}, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"user_id":"00000000-0000-0000-0000-000000000000","role":"user"}`,
		},
		{
			name:                 "Invalid_Header_Name",
			headerName:           "InvalidHeader",
			headerValue:          "Bearer valid_token",
			expectedStatusCode:   401,
			expectedResponseBody: `{"error":"empty auth header"}`,
		},
		{
			name:                 "Invalid_Header_Value",
			headerName:           "Authorization",
			headerValue:          "InvalidBearer valid_token",
			expectedStatusCode:   401,
			expectedResponseBody: `{"error":"invalid auth header"}`,
		},
		{
			name:                 "Empty_Token",
			headerName:           "Authorization",
			headerValue:          "",
			expectedStatusCode:   401,
			expectedResponseBody: `{"error":"token is empty"}`,
		},
		{
			name:        "Parse_Error",
			headerName:  "Authorization",
			headerValue: "Bearer invalid_token",
			token:       "invalid_token",
			mockBehavior: func(s *mock_service.MockUserService, token string) {
				s.EXPECT().ParseToken(token).Return(nil, errors.New("invalid token"))
			},
			expectedStatusCode:   401,
			expectedResponseBody: `{"error":"invalid token"}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Create a new instance of the mock controller
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			// Create a mock of the UserService
			mockUserService := mock_service.NewMockUserService(ctrl)
			test.mockBehavior(mockUserService, test.token)

			// Create an instance of the UserService implementation
			service := usecase.NewUserService(mockUserService)

			// Create an instance of the UserMiddleware with the UserService dependency
			middleware := middleware.NewUserMiddleware(service)
			// Create a test HTTP handler function
			handler := middleware.Authenticate(func(w http.ResponseWriter, r *http.Request) {
				// Your handler logic goes here
			})

			// Create a test request with the specified headers
			req := httptest.NewRequest("GET", "/api/v1/signIn", nil)
			req.Header.Set(test.headerName, test.headerValue)

			// Create a response recorder to capture the response
			w := httptest.NewRecorder()

			// Serve the request through the middleware and handler
			handler.ServeHTTP(w, req)

			// Check the response status code
			if w.Code != test.expectedStatusCode {
				t.Errorf("expected status code %d, got %d", test.expectedStatusCode, w.Code)
			}

			// Check the response body
			body := w.Body.String()
			if body != test.expectedResponseBody {
				t.Errorf("expected response body '%s', got '%s'", test.expectedResponseBody, body)
			}
		})
	}
}
