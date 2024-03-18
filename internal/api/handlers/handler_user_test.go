package handlers

import (
	"bytes"
	mock_service "cinema_service/internal/api/handlers/mocks"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

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
