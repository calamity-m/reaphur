package bindings

// This provides common context keys bindings
type RequestIDKey struct{}

// provides common global command bindings
var (
	Debug bool
)

// provides some nice defaults to rely on for sweet sweet enmeshment
const (
	DefaultCentralAddress = "127.0.0.1:9001"
	DefaultGWAddress      = "127.0.0.1:9002"
	DefaultRedisAddress   = "127.0.0.1:6379"
)
