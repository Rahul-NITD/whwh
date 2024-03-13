package drivers

import (
	"net/http"

	"github.com/aargeee/whwh/AcceptanceTesting/specs"
)

type ATDriver struct{}

// Dispatch implements specs.WebhookTesterSubject.
func (ATDriver) Dispatch(req *http.Request, hook_url string) (res *http.Response, err error) {
	return nil, nil
}

// ConnectClientAndServer implements specs.WebhookTesterSubject.
func (ATDriver) ConnectClientAndServer(chanID string) (client specs.Client, err error) {
	return nil, nil
}

// CreateChannel implements specs.WebhookTesterSubject.
func (ATDriver) CreateChannel() (chanID string, err error) {
	return "", nil
}
