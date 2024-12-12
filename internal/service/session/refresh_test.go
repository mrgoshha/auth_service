package session

import (
	"AuthenticationService/internal/model"
	"AuthenticationService/pkg/auth"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"

	repomocks "AuthenticationService/internal/adapter/dbs/mocks"
	servicemocks "AuthenticationService/internal/service/mocks"
	managermocks "AuthenticationService/pkg/auth/mocks"
	sendermocks "AuthenticationService/pkg/email_sender/mocks"
)

func TestService_Refresh(t *testing.T) {
	// Init Test Table
	type mockBehaviorRepo func(r *repomocks.MockRefreshSessionRepository, oldSessionId string)
	type mockBehaviorService func(r *servicemocks.MockUser, id string)
	type mockBehaviorTokenManager func(r *managermocks.MockTokenManager, accToken, userId, userIP, sessionID string)
	type mockBehaviorTEmailSender func(r *sendermocks.MockEmailSender, userEmail, msg string)

	tests := []struct {
		name            string
		inputIp         string
		inputId         string
		inputAccToken   string
		inputRefToken   string
		mbService       mockBehaviorService
		mbRepo          mockBehaviorRepo
		mbTokenManager  mockBehaviorTokenManager
		mbEmailSender   mockBehaviorTEmailSender
		expectedAccess  string
		expectedRefresh string
		expectedError   error
	}{
		{
			name:          "Ok",
			inputIp:       "127.0.0.1:2323",
			inputId:       "user_id",
			inputAccToken: "oldAccToken",
			inputRefToken: "841fecba963761f240a1c291398385f01e99ee941450e64c481a4ec9a24686c1",
			mbService: func(r *servicemocks.MockUser, id string) {
				user := &model.User{
					Id:    "user_id",
					Email: "email",
				}
				r.EXPECT().GetUserById(id).Return(user, nil)
			},
			mbRepo: func(r *repomocks.MockRefreshSessionRepository, oldSessionId string) {
				oldSession := &model.Session{
					Id:           1,
					SessionId:    oldSessionId,
					RefreshToken: "$2a$10$mqTKnGs/35Wj7wSsDlTI6.vII2m/7L9/EN.ctPFAkD6OGM15cN5s2",
					Ip:           "127.0.0.1:2323",
					UserId:       "user_id",
					ExpiresAt:    time.Now().Add(time.Hour),
				}
				r.EXPECT().GetSessionBySessionId(oldSessionId).Return(oldSession, nil)
				r.EXPECT().UpdateSession(gomock.AssignableToTypeOf(oldSession)).Return(nil)
			},
			mbTokenManager: func(r *managermocks.MockTokenManager, accToken, userId, userIP, sessionID string) {
				payload := &auth.TokenPayload{
					UserId:    "user_id",
					UserIP:    "ip",
					SessionId: "session_id",
				}
				r.EXPECT().Parse(accToken).Return(payload, nil)
				r.EXPECT().NewJWT(userId, userIP, gomock.Any()).Return("access", nil)
				r.EXPECT().NewRefreshToken().Return("refresh", nil)
			},
			mbEmailSender:   func(r *sendermocks.MockEmailSender, userEmail, msg string) {},
			expectedAccess:  "access",
			expectedRefresh: "refresh",
			expectedError:   nil,
		},
		{
			name:          "Change ip",
			inputIp:       "127.0.0.1:2333",
			inputId:       "user_id",
			inputAccToken: "oldAccToken",
			inputRefToken: "841fecba963761f240a1c291398385f01e99ee941450e64c481a4ec9a24686c1",
			mbService: func(r *servicemocks.MockUser, id string) {
				user := &model.User{
					Id:    "user_id",
					Email: "email",
				}
				r.EXPECT().GetUserById(id).Return(user, nil)
			},
			mbRepo: func(r *repomocks.MockRefreshSessionRepository, oldSessionId string) {
				oldSession := &model.Session{
					Id:           1,
					SessionId:    oldSessionId,
					RefreshToken: "$2a$10$mqTKnGs/35Wj7wSsDlTI6.vII2m/7L9/EN.ctPFAkD6OGM15cN5s2",
					Ip:           "127.0.0.1:2323",
					UserId:       "user_id",
					ExpiresAt:    time.Now().Add(time.Hour),
				}
				r.EXPECT().GetSessionBySessionId(oldSessionId).Return(oldSession, nil)
				r.EXPECT().UpdateSession(gomock.AssignableToTypeOf(oldSession)).Return(nil)
			},
			mbTokenManager: func(r *managermocks.MockTokenManager, accToken, userId, userIP, sessionID string) {
				payload := &auth.TokenPayload{
					UserId:    "user_id",
					UserIP:    "ip",
					SessionId: "session_id",
				}
				r.EXPECT().Parse(accToken).Return(payload, nil)
				r.EXPECT().NewJWT(userId, userIP, gomock.Any()).Return("access", nil)
				r.EXPECT().NewRefreshToken().Return("refresh", nil)
			},
			mbEmailSender: func(r *sendermocks.MockEmailSender, userEmail, msg string) {
				r.EXPECT().SendEmail(gomock.AssignableToTypeOf(userEmail), gomock.AssignableToTypeOf(msg)).Return(nil)
			},
			expectedAccess:  "access",
			expectedRefresh: "refresh",
			expectedError:   nil,
		},
		{
			name:          "Expired Token",
			inputIp:       "127.0.0.1:2323",
			inputId:       "user_id",
			inputAccToken: "oldAccToken",
			inputRefToken: "841fecba963761f240a1c291398385f01e99ee941450e64c481a4ec9a24686c1",
			mbService: func(r *servicemocks.MockUser, id string) {
				user := &model.User{
					Id:    "user_id",
					Email: "email",
				}
				r.EXPECT().GetUserById(id).Return(user, nil)
			},
			mbRepo: func(r *repomocks.MockRefreshSessionRepository, oldSessionId string) {
				oldSession := &model.Session{
					Id:           1,
					SessionId:    oldSessionId,
					RefreshToken: "$2a$10$mqTKnGs/35Wj7wSsDlTI6.vII2m/7L9/EN.ctPFAkD6OGM15cN5s2",
					Ip:           "127.0.0.1:2323",
					UserId:       "user_id",
					ExpiresAt:    time.Now(),
				}
				r.EXPECT().GetSessionBySessionId(oldSessionId).Return(oldSession, nil)
			},
			mbTokenManager: func(r *managermocks.MockTokenManager, accToken, userId, userIP, sessionID string) {
				payload := &auth.TokenPayload{
					UserId:    "user_id",
					UserIP:    "ip",
					SessionId: "session_id",
				}
				r.EXPECT().Parse(accToken).Return(payload, nil)
			},
			mbEmailSender:   func(r *sendermocks.MockEmailSender, userEmail, msg string) {},
			expectedAccess:  "",
			expectedRefresh: "",
			expectedError:   fmt.Errorf("token is expired"),
		},
		{
			name:          "Invalid Refresh Token",
			inputIp:       "127.0.0.1:2323",
			inputId:       "user_id",
			inputAccToken: "oldAccToken",
			inputRefToken: "841fecba963761f240a1c291398385f01e99ee941450e64c481a4ec9a24686c3",
			mbService: func(r *servicemocks.MockUser, id string) {
				user := &model.User{
					Id:    "user_id",
					Email: "email",
				}
				r.EXPECT().GetUserById(id).Return(user, nil)
			},
			mbRepo: func(r *repomocks.MockRefreshSessionRepository, oldSessionId string) {
				oldSession := &model.Session{
					Id:           1,
					SessionId:    oldSessionId,
					RefreshToken: "$2a$10$mqTKnGs/35Wj7wSsDlTI6.vII2m/7L9/EN.ctPFAkD6OGM15cN5s2",
					Ip:           "127.0.0.1:2323",
					UserId:       "user_id",
					ExpiresAt:    time.Now(),
				}
				r.EXPECT().GetSessionBySessionId(oldSessionId).Return(oldSession, nil)
			},
			mbTokenManager: func(r *managermocks.MockTokenManager, accToken, userId, userIP, sessionID string) {
				payload := &auth.TokenPayload{
					UserId:    "user_id",
					UserIP:    "ip",
					SessionId: "session_id",
				}
				r.EXPECT().Parse(accToken).Return(payload, nil)
			},
			mbEmailSender:   func(r *sendermocks.MockEmailSender, userEmail, msg string) {},
			expectedAccess:  "",
			expectedRefresh: "",
			expectedError:   fmt.Errorf("token is not valid crypto/bcrypt: hashedPassword is not the hash of the given password"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Init Dependencies
			c := gomock.NewController(t)
			defer c.Finish()

			tokenManager := managermocks.NewMockTokenManager(c)
			userService := servicemocks.NewMockUser(c)
			repo := repomocks.NewMockRefreshSessionRepository(c)
			sender := sendermocks.NewMockEmailSender(c)

			test.mbTokenManager(tokenManager, test.inputAccToken, test.inputId, test.inputIp, "")
			test.mbService(userService, test.inputId)
			test.mbRepo(repo, "session_id")
			test.mbEmailSender(sender, "", "")

			duration, _ := time.ParseDuration("2h")

			service := &Service{
				repository:      repo,
				userService:     userService,
				tokenManager:    tokenManager,
				refreshTokenTTL: duration,
				sender:          sender,
			}

			// Call method
			access, refresh, err := service.Refresh(test.inputAccToken, test.inputRefToken, test.inputIp)

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
