package systems_test

import (
	"testing"

	"github.com/Rahul-NITD/whwh/drivers"
	"github.com/Rahul-NITD/whwh/specs"
)

func TestSystem(t *testing.T) {
	specs.TesterSpecification(t, drivers.NewSysDriver())
	specs.TesterSpecification(t, drivers.DocDriver{})
}
