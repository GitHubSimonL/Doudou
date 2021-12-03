package _default

import "time"

const (
	DefaultPort             = 10086
	DefaultIP               = "127.0.0.1"
	DefaultPackageHeaderLen = 8
	DefaultRequestQueueLen  = 256
	DefaultConnectionTTL    = 1 * time.Minute
)
