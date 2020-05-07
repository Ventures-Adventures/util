package redis

import(
	"time"
)

type Options struct {
	Addr string
	Password string
	DB int
	DialTimeout time.Duration
	ReadTimeout time.Duration
	WriteTimeout time.Duration
	PoolSize int
	PoolTimeout time.Duration
}