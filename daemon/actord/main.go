package main

import (
	"fmt"
	"github.com/funkygao/dragon/actor"
	"github.com/funkygao/golib/locking"
	"github.com/funkygao/golib/server"
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

	server.SetupLogging(options.logFile, options.logLevel)

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

	server := server.NewServer("actord")
	server.LoadConfig(options.configFile)
	server.Launch()

	actor := actor.New(server)
	actor.ServeForever()

	shutdown()
}
