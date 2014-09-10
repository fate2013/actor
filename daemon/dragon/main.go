package main

import (
	"fmt"
	"github.com/funkygao/dragon/actor"
	"github.com/funkygao/dragon/server"
	"github.com/funkygao/golib/locking"
	"github.com/funkygao/golib/signal"
	"os"
	"runtime/debug"
	"syscall"
)

func init() {
	parseFlags()

	if options.showVersion {
		server.ShowVersionAndExit()
	}

	if options.lockFile != "" {
		if locking.InstanceLocked(options.lockFile) {
			fmt.Fprintf(os.Stderr, "Another dragon is running, exit...\n")
			os.Exit(1)
		}

		locking.LockInstance(options.lockFile)
	}

	signal.RegisterSignalHandler(syscall.SIGINT, func(sig os.Signal) {
		shutdown()
	})

}

func main() {
	defer func() {
		cleanup()

		if err := recover(); err != nil {
			fmt.Println(err)
			debug.PrintStack()

		}
	}()

	server.SetupLogging(options.logFile, options.logLevel)

	server := server.NewServer("dragon")
	server.LoadConfig(options.configFile)
	server.Launch()

	actor := actor.NewActor(server)
	actor.ServeForever()

	shutdown()
}
