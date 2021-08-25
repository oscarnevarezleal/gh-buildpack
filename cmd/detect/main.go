package main

import (
	"github.com/cloudfoundry/packit"
	"laraboot-buildpacks/gh/laraboot"
)

func main() {
	packit.Detect(laraboot.Detect())
}
