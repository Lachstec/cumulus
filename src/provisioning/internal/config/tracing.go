package config

// TracingConfig contains configuration for connecting to a Jaeger instance for
// telemetry and structured logging.
type TracingConfig struct {
	// Endpoint is the URL to the Jaeger instance to send spans to.
	Endpoint string
	// ServiceName is the name that will be used to identify this service.
	ServiceName string
	// InsecureConnection specifies whether to validate Certificates or not. Should only be set when developing.
	InsecureConnection bool
}
