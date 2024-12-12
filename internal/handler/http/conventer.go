package http

import (
	apimodel "AuthenticationService/internal/handler/http/model"
)

func ToTokenApiModel(access, refresh string) *apimodel.Tokens {
	return &apimodel.Tokens{
		AccessToken:  access,
		RefreshToken: refresh,
	}
}
