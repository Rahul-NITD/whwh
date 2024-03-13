package drivers

import (
	"net/http"
	"testing"

	"github.com/aargeee/whwh/AcceptanceTesting/specs"
)

type ATDriver struct{}

// ConnectClientAndServer implements specs.WebhookTesterSubject.
func (ATDriver) ConnectClientAndServer(t *testing.T, chanID string) (client specs.Client) {
	panic("unimplemented")
}

// CreateChannel implements specs.WebhookTesterSubject.
func (ATDriver) CreateChannel(t *testing.T) (chanID string) {
	panic("unimplemented")
}

// Dispatch implements specs.WebhookTesterSubject.
func (ATDriver) Dispatch(t *testing.T, req *http.Request, hook_url string) (res *http.Response) {
	panic("unimplemented")
}
