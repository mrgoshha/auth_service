package http

import (
	"fmt"
	"os"
	"time"
)

const (
	port         = "HTTP_PORT"
	host         = "HTTP_HOST"
	writeTimeout = "HTTP_WRITE_TIMEOUT"
	readTimeout  = "HTTP_READ_TIMEOUT"
)

type Config struct {
	Host         string
	Port         string
	WriteTimeout time.Duration
	ReadTimeout  time.Duration
}

func NewConfig() (*Config, error) {
	port := os.Getenv(port)
	if len(port) == 0 {
		return nil, fmt.Errorf("config error")
	}
	host := os.Getenv(host)
	if len(host) == 0 {
		return nil, fmt.Errorf("config error")
	}
	writeTimeout := os.Getenv(writeTimeout)
	if len(writeTimeout) == 0 {
		writeTimeout = "15s"
	}
	readTimeout := os.Getenv(readTimeout)
	if len(readTimeout) == 0 {
		readTimeout = "15s"
	}

	wt, err := time.ParseDuration(writeTimeout)
	if err != nil {
		return nil, fmt.Errorf("not duration format")
	}
	rt, err := time.ParseDuration(readTimeout)
	if err != nil {
		return nil, fmt.Errorf("not duration format")
	}
	return &Config{
		Host:         host,
		Port:         port,
		WriteTimeout: wt,
		ReadTimeout:  rt,
	}, nil
}
