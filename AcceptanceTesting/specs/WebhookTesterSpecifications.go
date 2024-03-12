package specs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type WebhookTesterSubject interface {
	CreateChannel() (chanID string, err error)
	ConnectClientAndServer(chanID string)
}

func TestWebhookTester(t *testing.T, subject WebhookTesterSubject) {
	_, err := subject.CreateChannel()
	assert.NoError(t, err, "Error while creating Channel")
}
