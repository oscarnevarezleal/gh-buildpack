# gh-buildpack

It provides GitHub’s official command line tool.

## Usage


## Build
```shell
pack buildpack package laraboot-buildpacks/gh --config ./package.toml
#test
pack build app --path . --buildpack laraboot-buildpacks/gh --builder paketobuildpacks/builder:tiny
```