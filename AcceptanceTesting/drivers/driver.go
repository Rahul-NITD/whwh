package drivers

type ATDriver struct{}

// ConnectClientAndServer implements specs.WebhookTesterSubject.
func (ATDriver) ConnectClientAndServer(chanID string) {
	panic("unimplemented")
}

// CreateChannel implements specs.WebhookTesterSubject.
func (ATDriver) CreateChannel() (chanID string, err error) {
	return "", nil
}
