package httpclient

import (
	"github.com/goinbox/golog"

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
	LogLevel int

	Timeout           time.Duration
	KeepAliveTime     time.Duration
	DisableKeepAlives bool

	MaxIdleConnsPerHost int
	MaxIdleConns        int
	IdleConnTimeout     time.Duration
}

func NewConfig() *Config {
	return &Config{
		LogLevel: golog.LevelDebug,

		Timeout:           DefaultTimeout,
		KeepAliveTime:     DefaultKeepaliveTime,
		DisableKeepAlives: false,

		MaxIdleConnsPerHost: DefaultMaxIdleConnsPerHost,
		MaxIdleConns:        DefaultMaxIdleConns,
		IdleConnTimeout:     DefaultIdleConnTimeout,
	}
}
