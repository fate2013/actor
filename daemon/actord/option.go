package main

import (
	"flag"
)

var (
	options struct {
		configFile  string
		showVersion bool
		logFile     string
		logLevel    string
		kill        bool
		lockFile    string
	}
)

func parseFlags() {
	flag.BoolVar(&options.kill, "k", false, "kill actord")
	flag.StringVar(&options.lockFile, "lockfile", "actord.lock", "lock file")
	flag.StringVar(&options.configFile, "conf", "etc/actord.cf", "config file")
	flag.StringVar(&options.logFile, "log", "stdout", "log file")
	flag.StringVar(&options.logLevel, "level", "debug", "log level")
	flag.BoolVar(&options.showVersion, "version", false, "show version and exit")

	flag.Parse()
}
