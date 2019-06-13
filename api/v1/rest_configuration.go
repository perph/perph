package v1

// RESTSpec defines the configuration options for a REST request
type RESTSpec struct {
	Path   string `json:"path"`
	Method string `json:"method"`

	// +optional
	QueryParams map[string]string `json:"queryParams,omitempty"`

	// +optional
	Headers map[string]string `json:"headers,omitempty"`

	// +optional
	Data []byte `json:"data,omitempty"`
}
