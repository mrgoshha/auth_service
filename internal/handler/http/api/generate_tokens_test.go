package api

import (
	httpApi "AuthenticationService/internal/handler/http"
	servicemocks "AuthenticationService/internal/service/mocks"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestAuthController_GenerateTokens(t *testing.T) {
	// Init Test Table
	type mockBehavior func(r *servicemocks.MockAuthSession, id, ip string)

	tests := []struct {
		name                 string
		inputId              string
		url                  string
		ip                   string
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:    "Ok",
			inputId: "43503be1-3a71-4135-a7ea-b42361957c56",
			url:     "/auth/token?id=43503be1-3a71-4135-a7ea-b42361957c56",
			ip:      "127.0.0.1:2323",
			mockBehavior: func(r *servicemocks.MockAuthSession, id, ip string) {
				r.EXPECT().GenerateTokens(id, ip).Return("access", "refresh", nil)
			},
			expectedStatusCode: http.StatusOK,
			expectedResponseBody: `{"access_token":"access","refresh_token":"refresh"}
`,
		},
		{
			name:               "Wrong Input",
			url:                "/auth/token?id=",
			mockBehavior:       func(r *servicemocks.MockAuthSession, id, ip string) {},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:    "Error in service",
			inputId: "43503be1-3a71-4135-a7ea-b42361957c56",
			url:     "/auth/token?id=43503be1-3a71-4135-a7ea-b42361957c56",
			ip:      "127.0.0.1:2323",
			mockBehavior: func(r *servicemocks.MockAuthSession, id, ip string) {
				r.EXPECT().GenerateTokens(id, ip).Return("", "", fmt.Errorf("smth error"))
			},
			expectedStatusCode: http.StatusInternalServerError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Init Dependencies
			c := gomock.NewController(t)
			defer c.Finish()

			service := servicemocks.NewMockAuthSession(c)
			test.mockBehavior(service, test.inputId, test.ip)

			log := slog.New(
				slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

			controller := &AuthController{
				service: service,
				log:     log,
			}

			// Init Endpoint
			r := httpApi.NewRouter(log)
			r.HandleFunc("/auth/token", controller.GenerateTokens).Methods(http.MethodGet)

			// Create Request
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, test.url, nil)
			req.RemoteAddr = test.ip

			// Make Request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, w.Code, test.expectedStatusCode)
			if w.Code == http.StatusOK {
				assert.Equal(t, w.Body.String(), test.expectedResponseBody)
			}
		})
	}
}
