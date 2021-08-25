package laraboot_test

import (
	"testing"

	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
)

func TestUnitGoBuild(t *testing.T) {
	suite := spec.New("go-build", spec.Report(report.Terminal{}))
	suite("Build", testBuild)
	suite.Run(t)
}
