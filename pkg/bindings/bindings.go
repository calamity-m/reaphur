package bindings

// This provides common context keys bindings
type RequestIDKey struct{}

// provides common global command bindings
var (
	Debug bool
)

// provides some nice defaults to rely on for sweet sweet enmeshment
const (
	DefaultCentralAddress = "localhost:9001"
	DefaultGWAddress      = "localhost:9002"
)
