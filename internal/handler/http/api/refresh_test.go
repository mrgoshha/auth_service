package api

import (
	httpApi "AuthenticationService/internal/handler/http"
	servicemocks "AuthenticationService/internal/service/mocks"
	"bytes"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestAuthController_Refresh(t *testing.T) {
	// Init Test Table
	type mockBehavior func(r *servicemocks.MockAuthSession, access, refresh, ip string)

	tests := []struct {
		name                 string
		inputBody            string
		access               string
		refresh              string
		ip                   string
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name: "Ok",
			inputBody: `{
				"access_token": "access",
				"refresh_token": "refresh"
			}`,
			access:  "access",
			refresh: "refresh",
			ip:      "127.0.0.1:2323",
			mockBehavior: func(r *servicemocks.MockAuthSession, access, refresh, ip string) {
				r.EXPECT().Refresh(access, refresh, ip).Return("new_access", "new_refresh", nil)
			},
			expectedStatusCode: http.StatusOK,
			expectedResponseBody: `{"access_token":"new_access","refresh_token":"new_refresh"}
`,
		},
		{
			name:               "Wrong Input",
			inputBody:          "Invalid",
			mockBehavior:       func(r *servicemocks.MockAuthSession, access, refresh, ip string) {},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name: "Error in service",
			inputBody: `{
				"access_token": "access",
				"refresh_token": "refresh"
			}`,
			access:  "access",
			refresh: "refresh",
			ip:      "127.0.0.1:2323",
			mockBehavior: func(r *servicemocks.MockAuthSession, access, refresh, ip string) {
				r.EXPECT().Refresh(access, refresh, ip).Return("", "", fmt.Errorf("smth error"))
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
			test.mockBehavior(service, test.access, test.refresh, test.ip)

			log := slog.New(
				slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

			controller := &AuthController{
				service: service,
				log:     log,
			}

			// Init Endpoint
			r := httpApi.NewRouter(log)
			r.HandleFunc("/auth/refresh", controller.Refresh).Methods(http.MethodPost)

			// Create Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/auth/refresh",
				bytes.NewBufferString(test.inputBody))
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
