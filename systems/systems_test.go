package systems_test

import (
	"testing"

	"github.com/Rahul-NITD/whwh/drivers"
	"github.com/Rahul-NITD/whwh/specs"
)

func TestSystem(t *testing.T) {
	driver := &drivers.SysDriver{}
	specs.TesterSpecification(t, driver)
}
