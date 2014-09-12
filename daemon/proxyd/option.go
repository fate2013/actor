package main

import (
	"flag"
)

func parseFlags() {
	flag.StringVar(&options.configFile, "conf", "etc/proxyd.cf", "config file")
	flag.StringVar(&options.logFile, "log", "stdout", "log file")
	flag.StringVar(&options.logLevel, "level", "debug", "log level")

	flag.Parse()
}
