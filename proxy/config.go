package proxy

import (
	conf "github.com/funkygao/jsconf"
	log "github.com/funkygao/log4go"
	"time"
)

type proxyConfig struct {
	statsInterval int
	tcpNoDelay    bool
	tcpIoTimeout  time.Duration
	proxyPass     string // host:port, TODO support server farm
	pm            pmConfig
}

// process management, naming after php-fpm
type pmConfig struct {
	maxServerN      int
	startServerN    int
	minSpareServerN int
	maxSpareServerN int
	spawnBatchSize  int
}

func (this *Proxy) loadConfig(cf *conf.Conf) {
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
	this.maxServerN = section.Int("max_servers", 10)
	this.startServerN = section.Int("start_servers", 5)
	this.minSpareServerN = section.Int("min_spare_servers", 3)
	this.maxSpareServerN = section.Int("max_spare_servers", this.minSpareServerN+2)
	this.spawnBatchSize = section.Int("spawn_batch_size", 3)
	if this.minSpareServerN > this.maxSpareServerN {
		panic("error: min_spare_servers > max_spare_servers")
	}
	if this.startServerN > this.maxServerN {
		panic("error: start_servers > max_servers")
	}
}
