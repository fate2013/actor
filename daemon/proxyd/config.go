package main

import (
	conf "github.com/funkygao/jsconf"
	log "github.com/funkygao/log4go"
	"time"
)

type proxyConfig struct {
	statsInterval int
	tcpNoDelay    bool
	tcpIoTimeout  time.Duration
	proxyPass     string // host:port
	pm            pmConfig
}

// process management, naming from php-fpm
type pmConfig struct {
	maxOutstandingSessionN int // a session is a persistent conn with upstream
	startServerN           int
	minSpareServerN        int
	spawnBatchSize         int
}

func (this *proxy) loadConfig(cf *conf.Conf) {
	this.config = proxyConfig{}
	this.config.statsInterval = cf.Int("stats_interval", 5)
	this.config.tcpNoDelay = cf.Bool("tcp_nodelay", true)
	this.config.tcpIoTimeout = time.Duration(cf.Int("tcp_io_timeout", 5)) * time.Second
	this.config.proxyPass = cf.StringList("proxy_pass", nil)[0] // TODO

	// pm section
	this.config.pm = pmConfig{}
	section, err := cf.Section("pm")
	if err != nil {
		panic(err)
	}
	this.config.pm.loadConfig(section)

	log.Debug("config loaded: %#v", this.config)
}

func (this *pmConfig) loadConfig(section *conf.Conf) {
	this.maxOutstandingSessionN = section.Int("max_outstanding_sessions", 10)
	this.startServerN = section.Int("start_servers", 5)
	this.minSpareServerN = section.Int("min_spare_servers", 3)
	this.spawnBatchSize = section.Int("spawn_batch_size", 3)
}
