package drivers

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/aargeee/whwh/AcceptanceTesting/specs"
	"github.com/aargeee/whwh/whwh"
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
	req, err := http.NewRequest(http.MethodPost, d.Server_url+"/create", http.NoBody)
	assert.NoError(t, err, "Could not build /create request")

	res, err := http.DefaultClient.Do(req)
	assert.NoError(t, err, "Could not execute /create request")

	var response whwh.CreateChannelResponse
	assert.NoError(t, json.NewDecoder(res.Body).Decode(&response), "Could not unmarshal response")

	assert.Equal(t, "CreateChannel", response.Event)
	assert.Equal(t, "SUCCESS", response.Status)
	assert.Equal(t, "Channel Created Successfully", response.Message)
	return response.Response.ChannelID
}

// Dispatch implements specs.WebhookTesterSubject.
func (ATDriver) Dispatch(t *testing.T, req *http.Request, hook_url string) (res *http.Response) {
	panic("unimplemented Dispatch")
}
