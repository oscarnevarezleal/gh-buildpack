package main

import (
	"github.com/paketo-buildpacks/packit"
	"github.com/paketo-buildpacks/packit/chronos"
	laraboot "laraboot-buildpacks/gh/laraboot"
	"log"
	"os"
)

func init() {
	log.Println("::init")
}

func main() {

	logEmitter := laraboot.NewLogEmitter(os.Stdout)

	packit.Build(laraboot.Build(
		logEmitter,
		chronos.DefaultClock))
}
