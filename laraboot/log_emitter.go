package laraboot

import (
	"io"
	"strconv"

	"github.com/paketo-buildpacks/packit"
	"github.com/paketo-buildpacks/packit/scribe"
)

type LogEmitter struct {
	// Emitter is embedded and therefore delegates all of its functions to the
	// LogEmitter.
	scribe.Emitter
}

func NewLogEmitter(output io.Writer) LogEmitter {
	return LogEmitter{
		Emitter: scribe.NewEmitter(output),
	}
}

func (e LogEmitter) Environment(environment packit.Environment) {
	e.Logger.Subprocess("%s", scribe.NewFormattedMapFromEnvironment(environment))
}

func (e LogEmitter) Candidates(entries []packit.BuildpackPlanEntry) {
	e.Subprocess("Candidate version sources (in priority order):")

	var (
		sources [][2]string
		maxLen  int
	)

	for _, entry := range entries {
		versionSource, ok := entry.Metadata["version-source"].(string)
		if !ok {
			versionSource = "<unknown>"
		}

		version, ok := entry.Metadata["version"].(string)
		if !ok {
			version = "*"
		}

		if len(versionSource) > maxLen {
			maxLen = len(versionSource)
		}

		sources = append(sources, [2]string{versionSource, version})
	}

	for _, source := range sources {
		e.Action(("%-" + strconv.Itoa(maxLen) + "s -> %q"), source[0], source[1])
	}

	e.Break()
}
