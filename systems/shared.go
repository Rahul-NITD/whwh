package systems

type HealthReport struct {
	Status  string `json:"status"`
	Event   string `json:"event"`
	Message string `json:"message"`
}
