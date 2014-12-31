package config

import (
	"github.com/funkygao/golib/ip"
	conf "github.com/funkygao/jsconf"
	log "github.com/funkygao/log4go"
	"strings"
	"time"
)

type ActorConfig struct {
	EtcdServers       []string
	EtcdSelfAddr      string
	HttpApiListenAddr string
	StatsListenAddr   string
	ProfListenAddr    string
	MetricsLogfile    string
	SchedulerBacklog  int

	ScheduleInterval     time.Duration
	ConsoleStatsInterval time.Duration

	Mysql     ConfigMysql
	Beanstalk ConfigBeanstalk
	Worker    ConfigWorker
	Pubnub    ConfigPubnub
}

func (this *ActorConfig) LoadConfig(cf *conf.Conf) (err error) {
	this.EtcdServers = cf.StringList("etcd_servers", nil)
	if len(this.EtcdServers) > 0 {
		this.EtcdSelfAddr = cf.String("etcd_self_addr", "")
		if strings.HasPrefix(this.EtcdSelfAddr, ":") {
			// automatically get local ip addr
			myIp := ip.LocalIpv4Addrs()[0]
			this.EtcdSelfAddr = myIp + this.EtcdSelfAddr
		}
	}
	this.HttpApiListenAddr = cf.String("http_api_listen_addr", ":9898")
	this.StatsListenAddr = cf.String("stats_listen_addr", "127.0.0.1:9010")
	this.ProfListenAddr = cf.String("prof_listen_addr", "")
	this.MetricsLogfile = cf.String("metrics_logfile", "")
	this.SchedulerBacklog = cf.Int("sched_backlog", 10000)

	this.ScheduleInterval = cf.Duration("sched_interval", time.Second)
	this.ConsoleStatsInterval = cf.Duration("stats_interval", time.Minute*10)

	this.Mysql = ConfigMysql{}
	var section *conf.Conf
	section, err = cf.Section("mysql")
	if err != nil {
		return
	}
	this.Mysql.loadConfig(section)

	this.Beanstalk = ConfigBeanstalk{}
	section, err = cf.Section("beanstalk")
	if err != nil {
		return err
	}
	this.Beanstalk.loadConfig(section)

	this.Worker = ConfigWorker{}
	section, err = cf.Section("worker")
	if err != nil {
		return
	}
	this.Worker.loadConfig(section)

	this.Pubnub = ConfigPubnub{}
	section, err = cf.Section("pubnub")
	if err != nil {
		return
	}
	this.Pubnub.loadConfig(section)

	log.Debug("actor config %+v", *this)
	return
}
