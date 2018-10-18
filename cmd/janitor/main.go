package main

import (
	"gopkg.in/src-d/go-cli.v0"
)

var (
	version string
	build   string
)

var app = cli.New("janitor", version, build, "Repository mainteinance tasks.")

func main() {
	app.RunMain()
}
