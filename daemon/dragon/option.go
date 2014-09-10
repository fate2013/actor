package main

import (
	"flag"
)

func parseFlags() {
	flag.StringVar(&options.lockFile, "lockfile", "dragon.lock", "lock file")
	flag.StringVar(&options.configFile, "conf", "etc/dragon.cf", "config file")
	flag.StringVar(&options.logFile, "log", "stdout", "log file")
	flag.StringVar(&options.logLevel, "level", "debug", "log level")
	flag.BoolVar(&options.showVersion, "version", false, "show version and exit")

	flag.Parse()
}
