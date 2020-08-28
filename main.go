package main

import (
	"flag"
	"fmt"
	"os"

	"doodod.com/doodod/shortme/conf"
	"doodod.com/doodod/shortme/web"
)

func main() {
	env := os.Getenv("ENV_TYPE")
	var cfgFile *string
	cfgFile = flag.String("c", "config.conf", "configuration file")

	version := flag.Bool("v", false, "Version")

	flag.Parse()

	if *version {
		fmt.Println(conf.Version)
		os.Exit(0)
	}

	// parse config
	conf.MustParseConfig(*cfgFile)

	// api
	web.Start()
}
