package config

import (
	"flag"
	"fmt"
)

var (
	ConfigFile string
)

func init() {
	// 1) define flags (command line arguments)
	cfgFl := flag.String("config", "config.yml", "path to config file yaml (yml)")

	// 2) parse them into variables
	flag.Parse()

	// 3) write to global variable
	ConfigFile = *cfgFl

	fmt.Printf("config file path: %v"+"\n", ConfigFile)
}
