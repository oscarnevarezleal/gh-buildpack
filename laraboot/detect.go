package laraboot

import (
	"github.com/cloudfoundry/packit"
)

func Detect() packit.DetectFunc {
	return func(context packit.DetectContext) (packit.DetectResult, error) {
		return packit.DetectResult{
			Plan: packit.BuildPlan{
				Provides: []packit.BuildPlanProvision{
					{Name: "gh"},
				},
				Requires: []packit.BuildPlanRequirement{
					{
						Name:     "gh",
						Metadata: map[string]string{},
					},
				},
			},
		}, nil
	}
}
