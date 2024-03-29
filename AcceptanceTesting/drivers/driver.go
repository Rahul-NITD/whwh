package drivers

import (
	"net/http"
	"testing"

	"github.com/aargeee/whwh/AcceptanceTesting/specs"
	"github.com/aargeee/whwh/whwh/client"
	"github.com/stretchr/testify/assert"
)

type ATDriver struct {
	Server_url string
}

// ConnectClientAndServer implements specs.WebhookTesterSubject.
func (ATDriver) ConnectClientAndServer(t *testing.T, chanID string) (client specs.Client) {
	panic("unimplemented ConnectClientAndServer")
}

// CreateChannel implements specs.WebhookTesterSubject.
func (d *ATDriver) CreateChannel(t *testing.T) (chanID string) {
	id, err := client.RequestChannel(d.Server_url + "/create")
	assert.NoError(t, err, "Could not create Channel")
	return id
}

// Dispatch implements specs.WebhookTesterSubject.
func (ATDriver) Dispatch(t *testing.T, req *http.Request, hook_url string) (res *http.Response) {
	panic("unimplemented Dispatch")
}
