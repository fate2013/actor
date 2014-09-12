package main

import (
	"github.com/funkygao/dragon/proxy"
	"github.com/funkygao/dragon/server"
)

func init() {
	parseFlags()

	server.SetupLogging(options.logFile, options.logLevel)
}

func main() {
	server := server.NewServer("proxyd")
	server.LoadConfig(options.configFile)
	server.Launch()

	proxy := proxy.New()
	proxy.Start(server).ServeForever()

}
