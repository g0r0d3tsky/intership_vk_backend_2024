package middleware

import (
	mock_service "cinema_service/internal/api/middleware/mocks"
	"cinema_service/internal/domain"
	"cinema_service/internal/usecase"
	"context"
	"errors"

	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"go.uber.org/mock/gomock"
)

func TestAuthenticateExecute(t *testing.T) {
	c := gomock.NewController(t)
	defer c.Finish()

	mockService := mock_service.NewMockUserService(c)

	expectedUserInfo := usecase.UserInfo{
		UserID: uuid.UUID{00000000 - 0000 - 0000 - 0000 - 000000000000},
		Role:   "user",
	}

	mockService.EXPECT().ParseToken("token123").Return(&expectedUserInfo, nil)

	middleware := &UserMiddleware{service: mockService}

	handlerCalled := false
	var capturedUserInfo *usecase.UserInfo
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handlerCalled = true
		capturedUserInfo = r.Context().Value(UserCtx).(*usecase.UserInfo)

		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Bearer token123")

	recorder := httptest.NewRecorder()

	middleware.Authenticate(handler).ServeHTTP(recorder, req)

	assert.True(t, handlerCalled, "Handler not called")
	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.NotNil(t, capturedUserInfo, "UserInfo not captured")
	assert.Equal(t, expectedUserInfo.UserID, capturedUserInfo.UserID)
	assert.Equal(t, expectedUserInfo.Role, capturedUserInfo.Role)
	assert.Equal(t, http.StatusOK, recorder.Code)

}

func TestAuthenticateHeader(t *testing.T) {
	dummyError := errors.New("dummy error")
	type mockBehavior func(r *mock_service.MockUserService, token string)
	u := &usecase.UserInfo{
		UserID: uuid.UUID{00000000 - 0000 - 0000 - 0000 - 000000000000},
		Role:   "user",
	}
	testTable := []struct {
		name                 string
		headerName           string
		headerValue          string
		token                string
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:                 "Invalid Header Name",
			headerName:           "",
			headerValue:          "Bearer token",
			token:                "token",
			mockBehavior:         func(r *mock_service.MockUserService, token string) {},
			expectedStatusCode:   401,
			expectedResponseBody: `{"error":"empty auth header"}`,
		},
		{
			name:                 "Invalid Header Value",
			headerName:           "Authorization",
			headerValue:          "Bearr token",
			token:                "token",
			mockBehavior:         func(r *mock_service.MockUserService, token string) {},
			expectedStatusCode:   401,
			expectedResponseBody: `{"error":"invalid auth header"}`,
		},
		{
			name:                 "Empty Token",
			headerName:           "Authorization",
			headerValue:          "Bearer ",
			token:                "token",
			mockBehavior:         func(r *mock_service.MockUserService, token string) {},
			expectedStatusCode:   401,
			expectedResponseBody: `{"error":"token is empty"}`,
		},
		{
			name:        "Parse Error",
			headerName:  "Authorization",
			headerValue: "Bearer token",
			token:       "token",
			mockBehavior: func(r *mock_service.MockUserService, token string) {
				r.EXPECT().ParseToken(token).Return(u, dummyError)
			},
			expectedStatusCode:   401,
			expectedResponseBody: `{"error":"dummy error"}`,
		},
	}

	for _, test := range testTable {
		t.Run(test.name, func(t *testing.T) {

			c := gomock.NewController(t)
			defer c.Finish()

			service := mock_service.NewMockUserService(c)
			test.mockBehavior(service, test.token)

			middleware := &UserMiddleware{service: service}

			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			})

			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			req.Header.Set(test.headerName, test.headerValue)

			recorder := httptest.NewRecorder()

			middleware.Authenticate(handler).ServeHTTP(recorder, req)

			// Asserts
			assert.Equal(t, recorder.Code, test.expectedStatusCode)
			assert.Equal(t, recorder.Body.String(), test.expectedResponseBody)
		})
	}
}

func TestRequireAdmin(t *testing.T) {
	adminUser := &usecase.UserInfo{
		UserID: uuid.UUID{},
		Role:   domain.ADMIN,
	}
	User := &usecase.UserInfo{
		UserID: uuid.UUID{},
		Role:   domain.USER,
	}
	testTable := []struct {
		name                 string
		ctx                  context.Context
		expectedResponseBody string
	}{
		{
			name:                 "Correct Admin context",
			ctx:                  context.WithValue(context.Background(), UserCtx, adminUser),
			expectedResponseBody: "",
		},
		{
			name:                 "Correct User context",
			ctx:                  context.WithValue(context.Background(), UserCtx, User),
			expectedResponseBody: `{"error":"Access denied"}`,
		},
		{
			name:                 "Empty context",
			ctx:                  context.Background(),
			expectedResponseBody: `{"error":"User context not found"}`,
		},

	}

	for _, tc := range testTable {

		fakeHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})

		middleware := &UserMiddleware{}

		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		req = req.WithContext(tc.ctx)

		recorder := httptest.NewRecorder()

		middleware.RequireAdmin(fakeHandler).ServeHTTP(recorder, req)

			assert.Equal(t, recorder.Body.String(), tc.expectedResponseBody)
	
	}
}
