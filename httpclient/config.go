package httpclient

import (
	"github.com/goinbox/golog"

	"time"
)

const (
	DEFAULT_TIMEOUT        = 30 * time.Second
	DEFAULT_KEEPALIVE_TIME = 30 * time.Second

	DEFAULT_MAX_IDLE_CONNS_PER_HOST = 10
	DEFAULT_MAX_IDLE_CONNS          = 100
	DEFAULT_IDLE_CONN_TIMEOUT       = 30 * time.Second
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

		Timeout:           DEFAULT_TIMEOUT,
		KeepAliveTime:     DEFAULT_KEEPALIVE_TIME,
		DisableKeepAlives: false,

		MaxIdleConnsPerHost: DEFAULT_MAX_IDLE_CONNS_PER_HOST,
		MaxIdleConns:        DEFAULT_MAX_IDLE_CONNS,
		IdleConnTimeout:     DEFAULT_IDLE_CONN_TIMEOUT,
	}
}
