package auth

import (
	"fmt"
	"os"
	"time"
)

const (
	SigningKey      = "JWT_SIGNING_KEY"
	AccessTokenTTL  = "ACCESS_TOKEN_TTL"
	RefreshTokenTTL = "REFRESH_TOKEN_TTL"
)

type Config struct {
	SigningKey      string
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
}

func NewConfig() (*Config, error) {
	key := os.Getenv(SigningKey)
	if len(key) == 0 {
		return nil, fmt.Errorf("empty signing key")
	}
	aTtl := os.Getenv(AccessTokenTTL)
	if len(aTtl) == 0 {
		aTtl = "2h"
	}
	rTtl := os.Getenv(RefreshTokenTTL)
	if len(rTtl) == 0 {
		rTtl = "720h" // 30 days
	}

	at, err := time.ParseDuration(aTtl)
	if err != nil {
		return nil, fmt.Errorf("not duration format")
	}
	rt, _ := time.ParseDuration(rTtl)
	if err != nil {
		return nil, fmt.Errorf("not duration format")
	}
	return &Config{
		SigningKey:      SigningKey,
		AccessTokenTTL:  at,
		RefreshTokenTTL: rt,
	}, nil
}
