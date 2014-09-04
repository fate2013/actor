package main

import (
	"flag"
	log "github.com/funkygao/log4go"
)

func parseFlags() {
	flag.StringVar(&options.lockFile, "lockfile", "dragon.lock", "lock file")
	flag.StringVar(&options.configFile, "conf", "etc/dragon.cf", "config file")
	flag.StringVar(&options.logFile, "log", "stdout", "log file")
	flag.StringVar(&options.logLevel, "level", "debug", "log level")
	flag.BoolVar(&options.showVersion, "version", false, "show version and exit")

	flag.Parse()
}

func setupLogging(logFile, logLevel string) {
	level := log.DEBUG
	switch logLevel {
	case "info":
		level = log.INFO

	case "warn":
		level = log.WARNING

	case "error":
		level = log.ERROR
	}

	for _, filter := range log.Global {
		filter.Level = level
	}

	if logFile == "stdout" {
		log.AddFilter("stdout", level, log.NewConsoleLogWriter())
	} else {
		writer := log.NewFileLogWriter(logFile, false)
		log.AddFilter("file", level, writer)
		writer.SetFormat("[%d %T] [%L] (%S) %M")
		writer.SetRotate(true)
		writer.SetRotateSize(0)
		writer.SetRotateLines(0)
		writer.SetRotateDaily(true)

	}

}
