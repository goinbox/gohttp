package httpclient

import (
	"time"
)

const (
	DefaultTimeout       = 30 * time.Second
	DefaultKeepaliveTime = 30 * time.Second

	DefaultMaxIdleConnsPerHost = 10
	DefaultMaxIdleConns        = 100
	DefaultIdleConnTimeout     = 30 * time.Second
)

type Config struct {
	Timeout           time.Duration
	KeepAliveTime     time.Duration
	DisableKeepAlives bool

	MaxIdleConnsPerHost int
	MaxIdleConns        int
	IdleConnTimeout     time.Duration

	LogRequestBody  bool
	LogResponseBody bool
}

func NewConfig() *Config {
	return &Config{
		Timeout:           DefaultTimeout,
		KeepAliveTime:     DefaultKeepaliveTime,
		DisableKeepAlives: false,

		MaxIdleConnsPerHost: DefaultMaxIdleConnsPerHost,
		MaxIdleConns:        DefaultMaxIdleConns,
		IdleConnTimeout:     DefaultIdleConnTimeout,

		LogRequestBody:  true,
		LogResponseBody: true,
	}
}
