package core

// OperationResult is a generic operation envelope for kernel actions.
type OperationResult[T any] struct {
	Level       CapabilityLevel    `json:"level"`
	Success     bool               `json:"success"`
	Data        *T                 `json:"data,omitempty"`
	RawStatus   int                `json:"rawStatus"`
	RawBody     string             `json:"rawBody,omitempty"`
	RawHeaders  map[string][]string `json:"rawHeaders,omitempty"`
	Message     string             `json:"message,omitempty"`
	Diagnostics map[string]any     `json:"diagnostics,omitempty"`
}
