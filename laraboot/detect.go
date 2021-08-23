package laraboot

import (
	"encoding/json"
	"fmt"
	"github.com/cloudfoundry/packit"
	"os"
	"path/filepath"
)

func Detect() packit.DetectFunc {
	return func(context packit.DetectContext) (packit.DetectResult, error) {
		return packit.DetectResult{
			Plan: packit.BuildPlan{
				Provides: []packit.BuildPlanProvision{
					{Name: "gh"},
				},
			},
		}, nil
	}
}
