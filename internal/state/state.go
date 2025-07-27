package state

// HTTPHeaders defines the structure for HTTP headers
type HTTPHeaders struct {
	RequestId    string
	PlatformType string
	Platform     string
	Version      string
}

// NewHTTPHeaders creates a new HTTPHeaders instance
func NewHTTPHeaders() *HTTPHeaders {
	return &HTTPHeaders{
		RequestId:    "request_id",
		PlatformType: "platform_type",
		Platform:     "platform",
		Version:      "version",
	}
}

// Global instance
var httpHeaders = NewHTTPHeaders()

// HttpHeaders returns the global HTTP headers instance
func HttpHeaders() *HTTPHeaders {
	return httpHeaders
}
