package service


//go:generate mockgen -source=auth_service.go -destination=mocks/auth_service.go

type AuthSession interface {
	GenerateTokens(id, ip string) (string, string, error)
	Refresh(accessToken, refreshToken, ip string) (string, string, error)
}
