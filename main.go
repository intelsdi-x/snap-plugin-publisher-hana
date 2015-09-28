package main

import (
    "os"

    "github.com/intelsdi-x/pulse/control/plugin"
    "github.com/intelsdi-x/pulse/plugin/publisher/pulse-publisher-hana/hana"
)

func main() {
    meta := hana.Meta()
    plugin.Start(meta, hana.NewHANAPublisher(), os.Args[1])
}

