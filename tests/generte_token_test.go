package tests

import (
	"net/http"
	"net/http/httptest"
)

func (s *APITestSuite) TestGenerateToken_Success() {
	r := s.Require()

	// Create Request
	req, _ := http.NewRequest(http.MethodGet, "/auth/token?id=92503be1-3a71-4135-a7ea-b42361957c56", nil)

	// Make Request
	resp := httptest.NewRecorder()
	s.serviceProvider.HttpRouter().ServeHTTP(resp, req)

	//Assert
	r.Equal(http.StatusOK, resp.Result().StatusCode)

}

func (s *APITestSuite) TestGenerateToken_InvalidUser() {
	r := s.Require()

	//Arrange
	output := `get user by id record not found`

	// Create Request
	req, _ := http.NewRequest(http.MethodGet, "/auth/token?id=user", nil)

	// Make Request
	resp := httptest.NewRecorder()
	s.serviceProvider.HttpRouter().ServeHTTP(resp, req)

	//Assert
	r.Equal(http.StatusNotFound, resp.Result().StatusCode)
	r.Contains(resp.Body.String(), output)

}

func (s *APITestSuite) TestGenerateToken_InvalidUserID() {
	r := s.Require()

	//Arrange
	outputTokens := `invalid id`

	// Create Request
	req, _ := http.NewRequest(http.MethodGet, "/auth/token", nil)

	// Make Request
	resp := httptest.NewRecorder()
	s.serviceProvider.HttpRouter().ServeHTTP(resp, req)

	//Assert
	r.Equal(http.StatusBadRequest, resp.Result().StatusCode)
	r.Contains(resp.Body.String(), outputTokens)

}
