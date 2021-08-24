package laraboot

import (
	_ "embed"
	"errors"
	"fmt"
	"github.com/bitfield/script"
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

		duration, ghGenErr := clock.Measure(func() error {
			return dependencyService.Deliver(dependency, context.CNBPath, binPath, "/platform")
		})

		if ghGenErr != nil {
			return packit.BuildResult{}, ghGenErr
		}
		logger.Action("Completed in %s", duration.Round(time.Millisecond))
		logger.Break()

		logger.Process("Set up environment")
		thisLayer.SharedEnv.Append("PATH", binPath, ":")
		ghBin := fmt.Sprintf("%s/gh_%s_linux_arm64/bin", binPath, dependency.Version)
		thisLayer.SharedEnv.Default("GH_BIN", ghBin)
		logger.Environment(thisLayer.SharedEnv)
		// expanding path and setting GH_BIN for runtime use
		envErr := os.Setenv("PATH", fmt.Sprintf("%s:%s", os.Getenv("PATH"), ghBin))
		if envErr != nil {
			return packit.BuildResult{}, ghGenErr
		}
		envErr = os.Setenv("GH_BIN", ghBin)
		if envErr != nil {
			return packit.BuildResult{}, ghGenErr
		}
		logger.Action("Completed in %s", duration.Round(time.Millisecond))
		logger.Break()

		logger.Process("Checking installation")

		p := script.Exec(fmt.Sprintf("%s", "gh"))
		output, _ := p.String()
		fmt.Println(output)

		var exit int = p.ExitStatus()
		if exit != 0 {
			err1 := errors.New("Instalation check failed: command exited with a non-zero status")
			return packit.BuildResult{}, err1
		}

		return packit.BuildResult{
			Layers: []packit.Layer{thisLayer},
		}, nil
	}
}
