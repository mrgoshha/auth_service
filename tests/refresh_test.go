package tests

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
)

func (s *APITestSuite) TestRefresh_Success() {
	r := s.Require()

	//Arrange
	accessToken, _ := s.tokenManager.NewJWT("92503be1-3a71-4135-a7ea-b42361957c56", "ip", "5b0f8c90-764b-4971-b82a-48f9f4c71705")
	refreshToken := "80e1dfafe52db03413313d86fa65bd50cb4cc80238cbea19b697b32e766258bb"
	inputTokens := fmt.Sprintf(`{"access_token":"%s","refresh_token":"%s"}`, accessToken, refreshToken)

	// Create Request
	req := httptest.NewRequest("POST", "/auth/refresh",
		bytes.NewBufferString(inputTokens))

	// Make Request
	resp := httptest.NewRecorder()
	s.serviceProvider.HttpRouter().ServeHTTP(resp, req)

	//Assert
	r.Equal(http.StatusOK, resp.Result().StatusCode)
}

func (s *APITestSuite) TestRefresh_InvalidInput() {
	r := s.Require()

	s.serviceProvider.AuthService()

	//Arrange
	inputTokens := "invalid"

	// Create Request
	req := httptest.NewRequest("POST", "/auth/refresh",
		bytes.NewBufferString(inputTokens))

	// Make Request
	resp := httptest.NewRecorder()
	s.serviceProvider.HttpRouter().ServeHTTP(resp, req)

	//Assert
	r.Equal(http.StatusBadRequest, resp.Result().StatusCode)
}

func (s *APITestSuite) TestRefresh_InvalidAccess() {
	r := s.Require()

	//Arrange
	accessToken := "access"
	refreshToken := "80e1dfafe52db03413313d86fa65bd50cb4cc80238cbea19b697b32e766258bb"
	inputTokens := fmt.Sprintf(`{"access_token":"%s","refresh_token":"%s"}`, accessToken, refreshToken)

	// Create Request
	req := httptest.NewRequest("POST", "/auth/refresh",
		bytes.NewBufferString(inputTokens))

	// Make Request
	resp := httptest.NewRecorder()
	s.serviceProvider.HttpRouter().ServeHTTP(resp, req)

	//Assert
	r.Equal(http.StatusInternalServerError, resp.Result().StatusCode)
}

func (s *APITestSuite) TestRefresh_UnpairedTokens() {
	r := s.Require()

	//Arrange
	accessToken, _ := s.tokenManager.NewJWT("92503be1-3a71-4135-a7ea-b42361957c56", "ip", "another_session_id")
	refreshToken := "80e1dfafe52db03413313d86fa65bd50cb4cc80238cbea19b697b32e766258bb"
	inputTokens := fmt.Sprintf(`{"access_token":"%s","refresh_token":"%s"}`, accessToken, refreshToken)

	// Create Request
	req := httptest.NewRequest("POST", "/auth/refresh",
		bytes.NewBufferString(inputTokens))

	// Make Request
	resp := httptest.NewRecorder()
	s.serviceProvider.HttpRouter().ServeHTTP(resp, req)

	//Assert
	r.Equal(http.StatusNotFound, resp.Result().StatusCode)
}

func (s *APITestSuite) TestRefresh_InvalidRefresh() {
	r := s.Require()

	//Arrange
	accessToken, _ := s.tokenManager.NewJWT("92503be1-3a71-4135-a7ea-b42361957c56", "ip", "5b0f8c90-764b-4971-b82a-48f9f4c71705")
	refreshToken := "invalid"
	inputTokens := fmt.Sprintf(`{"access_token":"%s","refresh_token":"%s"}`, accessToken, refreshToken)

	// Create Request
	req := httptest.NewRequest("POST", "/auth/refresh",
		bytes.NewBufferString(inputTokens))

	// Make Request
	resp := httptest.NewRecorder()
	s.serviceProvider.HttpRouter().ServeHTTP(resp, req)

	//Assert
	r.Equal(http.StatusInternalServerError, resp.Result().StatusCode)
}
