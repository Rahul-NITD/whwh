package whwh_test

import (
	"net/http/httptest"
	"testing"

	"github.com/aargeee/whwh/AcceptanceTesting/drivers"
	"github.com/aargeee/whwh/AcceptanceTesting/specs"
	"github.com/aargeee/whwh/handlers"
)

func TestWHWH(t *testing.T) {
	svr := httptest.NewServer(handlers.NewServer())
	specs.TestWebhookTester(t, &drivers.ATDriver{
		Server_url: svr.URL,
	})
}
