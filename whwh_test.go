package whwh_test

import (
	"testing"

	"github.com/aargeee/whwh/AcceptanceTesting/drivers"
	"github.com/aargeee/whwh/AcceptanceTesting/specs"
)

func TestWHWH(t *testing.T) {
	specs.TestWebhookTester(t, drivers.ATDriver{})
}
