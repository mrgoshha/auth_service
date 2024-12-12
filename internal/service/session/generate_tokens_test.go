package session

import (
	"AuthenticationService/internal/model"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"

	repomocks "AuthenticationService/internal/adapter/dbs/mocks"
	servicemocks "AuthenticationService/internal/service/mocks"
	managermocks "AuthenticationService/pkg/auth/mocks"
)

func TestService_GenerateTokens(t *testing.T) {
	// Init Test Table
	type mockBehaviorRepo func(r *repomocks.MockRefreshSessionRepository, session *model.Session)
	type mockBehaviorService func(r *servicemocks.MockUser, id string)
	type mockBehaviorTokenManager func(r *managermocks.MockTokenManager, userId, userIP, sessionID string)

	tests := []struct {
		name            string
		inputIp         string
		inputId         string
		mbService       mockBehaviorService
		mbRepo          mockBehaviorRepo
		mbTokenManager  mockBehaviorTokenManager
		expectedAccess  string
		expectedRefresh string
		expectedError   error
	}{
		{
			name:    "Ok",
			inputId: "43503be1-3a71-4135-a7ea-b42361957c56",
			inputIp: "127.0.0.1:2323",
			mbService: func(r *servicemocks.MockUser, id string) {
				r.EXPECT().GetUserById(id).Return(nil, nil)
			},
			mbRepo: func(r *repomocks.MockRefreshSessionRepository, session *model.Session) {
				r.EXPECT().CreateSession(gomock.AssignableToTypeOf(session)).Return(nil)
			},
			mbTokenManager: func(r *managermocks.MockTokenManager, userId, userIP, sessionID string) {
				r.EXPECT().NewJWT(userId, userIP, gomock.Any()).Return("access", nil)
				r.EXPECT().NewRefreshToken().Return("refresh", nil)
			},
			expectedAccess:  "access",
			expectedRefresh: "refresh",
			expectedError:   nil,
		},
		{
			name:    "User service error",
			inputId: "43503be1-3a71-4135-a7ea-b42361957c56",
			mbService: func(r *servicemocks.MockUser, id string) {
				r.EXPECT().GetUserById(id).Return(nil, fmt.Errorf("smth error"))
			},
			mbRepo:          func(r *repomocks.MockRefreshSessionRepository, session *model.Session) {},
			mbTokenManager:  func(r *managermocks.MockTokenManager, userId, userIP, sessionID string) {},
			expectedAccess:  "",
			expectedRefresh: "",
			expectedError:   fmt.Errorf("smth error"),
		},
		{
			name:    "Token manager error",
			inputId: "43503be1-3a71-4135-a7ea-b42361957c56",
			inputIp: "127.0.0.1:2323",
			mbService: func(r *servicemocks.MockUser, id string) {
				r.EXPECT().GetUserById(id).Return(nil, nil)
			},
			mbRepo: func(r *repomocks.MockRefreshSessionRepository, session *model.Session) {},
			mbTokenManager: func(r *managermocks.MockTokenManager, userId, userIP, sessionID string) {
				r.EXPECT().NewJWT(userId, userIP, gomock.Any()).Return("", fmt.Errorf("smth error"))
			},
			expectedAccess:  "",
			expectedRefresh: "",
			expectedError:   fmt.Errorf("create access token smth error"),
		},
		{
			name:    "Repository error",
			inputId: "43503be1-3a71-4135-a7ea-b42361957c56",
			inputIp: "127.0.0.1:2323",
			mbService: func(r *servicemocks.MockUser, id string) {
				r.EXPECT().GetUserById(id).Return(nil, nil)
			},
			mbRepo: func(r *repomocks.MockRefreshSessionRepository, session *model.Session) {
				r.EXPECT().CreateSession(gomock.AssignableToTypeOf(session)).Return(fmt.Errorf("smth error"))
			},
			mbTokenManager: func(r *managermocks.MockTokenManager, userId, userIP, sessionID string) {
				r.EXPECT().NewJWT(userId, userIP, gomock.Any()).Return("access", nil)
				r.EXPECT().NewRefreshToken().Return("refresh", nil)
			},
			expectedAccess:  "",
			expectedRefresh: "",
			expectedError:   fmt.Errorf("smth error"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Init Dependencies
			c := gomock.NewController(t)
			defer c.Finish()

			repo := repomocks.NewMockRefreshSessionRepository(c)
			userService := servicemocks.NewMockUser(c)
			tokenManager := managermocks.NewMockTokenManager(c)

			test.mbService(userService, test.inputId)
			test.mbRepo(repo, nil)
			test.mbTokenManager(tokenManager, test.inputId, test.inputIp, "")

			duration, _ := time.ParseDuration("2h")

			service := &Service{
				repository:      repo,
				userService:     userService,
				tokenManager:    tokenManager,
				refreshTokenTTL: duration,
			}

			// Call method
			access, refresh, err := service.GenerateTokens(test.inputId, test.inputIp)

			// Assert
			assert.Equal(t, test.expectedAccess, access)
			assert.Equal(t, test.expectedRefresh, refresh)
			if test.expectedError == nil {
				assert.Equal(t, test.expectedError, err)
			} else {
				assert.Equal(t, test.expectedError.Error(), err.Error())
			}

		})
	}
}
