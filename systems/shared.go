package systems

type HealthReport struct {
	Status  string `json:"status"`
	Event   string `json:"event"`
	Message string `json:"message"`
}

type StreamPayloadResponse struct {
	Event   string        `json:"event"`
	Payload StreamPayload `json:"data"`
}

type StreamPayload struct {
	StreamID string `json:"stream_id"`
}
