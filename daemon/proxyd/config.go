package main

import (
	conf "github.com/funkygao/jsconf"
)

type proxyConfig struct {
	tcpNoDelay   bool
	dragonServer string // host:port
	pm           *pmConfig
}

type pmConfig struct {
	maxOutstandingSessionN int
	startServerN           int
	minSpareServerN        int
	spawnBatchSize         int
}

func (this *dragonPool) loadConfig(cf *conf.Conf) {
	this.config = proxyConfig{}
	this.config.tcpNoDelay = cf.Bool("tcp_nodelay", true)
	this.config.dragonServer = cf.StringList("dragons", nil)[0]

	// pm section
	this.config.pm = new(pmConfig)
	section, err := cf.Section("pm")
	if err != nil {
		panic(err)
	}
	this.config.pm.loadConfig(section)

}

func (this *pmConfig) loadConfig(section *conf.Conf) {
	this.maxOutstandingSessionN = section.Int("max_outstanding_sessions", 10)
	this.startServerN = section.Int("start_servers", 5)
	this.minSpareServerN = section.Int("min_spare_servers", 3)
	this.spawnBatchSize = section.Int("spawn_batch_size", 3)
}
