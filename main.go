package main

import (
	"os"

	"github.com/intelsdi-x/pulse-plugin-publisher-hana/hana"
	"github.com/intelsdi-x/pulse/control/plugin"
)

func main() {
	meta := hana.Meta()
	plugin.Start(meta, hana.NewHANAPublisher(), os.Args[1])
}
