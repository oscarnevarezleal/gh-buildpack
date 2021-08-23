package laraboot

import (
	_ "embed"
	"fmt"
	"github.com/paketo-buildpacks/packit"
	"github.com/paketo-buildpacks/packit/cargo"
	"github.com/paketo-buildpacks/packit/chronos"
	"github.com/paketo-buildpacks/packit/postal"
	"os"
	"path/filepath"
	"time"
)

func Build(logger LogEmitter, clock chronos.Clock) packit.BuildFunc {
	return func(context packit.BuildContext) (packit.BuildResult, error) {

		logger.Title("%s %s", context.BuildpackInfo.Name, context.BuildpackInfo.Version)

		thisLayer, _ := context.Layers.Get("gh")

		transport := cargo.NewTransport()
		dependencyService := postal.NewService(transport)

		dependency, dependencyErr := dependencyService.Resolve(filepath.Join(context.CNBPath, "buildpack.toml"), "gh", "default", context.Stack)
		if dependencyErr != nil {
			return packit.BuildResult{}, dependencyErr
		}
		binPath := fmt.Sprintf("%s/bin", thisLayer.Path)
		logger.Subprocess("Installing Gh %s %s into %s", dependency.Version, dependency.SHA256, binPath)

		duration, blueprintGenErr := clock.Measure(func() error {
			return dependencyService.Deliver(dependency, context.CNBPath, binPath, "/platform")
		})

		if blueprintGenErr != nil {
			return packit.BuildResult{}, blueprintGenErr
		}
		logger.Action("Completed in %s", duration.Round(time.Millisecond))
		logger.Break()

		logger.Process("Configuring environment")
		thisLayer.SharedEnv.Append("PATH", binPath, ":")
		blueprintBin := fmt.Sprintf("%s/gh", binPath)
		thisLayer.SharedEnv.Default("GH_BIN", blueprintBin)
		logger.Environment(thisLayer.SharedEnv)

		// expanding path and setting GH_BIN for runtime use
		envErr := os.Setenv("PATH", fmt.Sprintf("%s:%s", os.Getenv("PATH"), binPath))
		if envErr != nil {
			return packit.BuildResult{}, blueprintGenErr
		}
		envErr = os.Setenv("GH_BIN", blueprintBin)
		if envErr != nil {
			return packit.BuildResult{}, blueprintGenErr
		}

		return packit.BuildResult{
			Layers: []packit.Layer{thisLayer},
		}, nil
	}
}
