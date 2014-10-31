package actor

import (
	"fmt"
	conf "github.com/funkygao/jsconf"
	log "github.com/funkygao/log4go"
	"time"
)

type ActorConfig struct {
	ApiListenAddr    string
	StatsListenAddr  string
	ProfListenAddr   string
	MetricsLogfile   string
	SchedulerBacklog int

	ScheduleInterval     time.Duration
	ConsoleStatsInterval time.Duration

	MysqlConfig  *ConfigMysql
	WorkerConfig *ConfigWorker
}

func (this *ActorConfig) LoadConfig(cf *conf.Conf) (err error) {
	this.ApiListenAddr = cf.String("api_listen_addr", ":9898")
	this.StatsListenAddr = cf.String("stats_listen_addr", "127.0.0.1:9010")
	this.ProfListenAddr = cf.String("prof_listen_addr", "")
	this.MetricsLogfile = cf.String("metrics_logfile", "")
	this.SchedulerBacklog = cf.Int("sched_backlog", 10000)

	this.ScheduleInterval = cf.Duration("sched_interval", time.Second)
	this.ConsoleStatsInterval = cf.Duration("stats_interval", time.Minute*10)

	this.MysqlConfig = new(ConfigMysql)
	var section *conf.Conf
	section, err = cf.Section("mysql")
	if err != nil {
		return
	}
	this.MysqlConfig.loadConfig(section)

	this.WorkerConfig = new(ConfigWorker)
	section, err = cf.Section("worker")
	if err != nil {
		return
	}
	this.WorkerConfig.loadConfig(section)

	log.Debug("actor config %+v", *this)
	return
}

type ConfigWorker struct {
	DryRun           bool
	DebugLocking     bool
	Timeout          time.Duration
	MaxFlightEntries int
	LockExpires      time.Duration

	// if use php as worker, it's callback url template
	Job   string
	March string
	Pve   string
}

func (this *ConfigWorker) loadConfig(cf *conf.Conf) {
	this.DryRun = cf.Bool("dry_run", true)
	this.DebugLocking = cf.Bool("debug_locking", false)
	this.Timeout = time.Duration(cf.Int("timeout", 5)) * time.Second
	this.MaxFlightEntries = cf.Int("max_flight_entries", 100000)
	this.LockExpires = cf.Duration("lock_expires", time.Second*30)
	this.Job = cf.String("job", "")
	this.March = cf.String("march", "")
	this.Pve = cf.String("pve", "")

	log.Debug("worker config: %+v", *this)
}

type ConfigMysql struct {
	ConnectTimeout time.Duration // part of DSN
	SlowThreshold  time.Duration

	Query   ConfigMysqlQuery
	Breaker ConfigBreaker

	Servers map[string]*ConfigMysqlInstance // key is pool
}

func (this *ConfigMysql) loadConfig(cf *conf.Conf) {
	this.ConnectTimeout = cf.Duration("connect_timeout", 4*time.Second)
	this.SlowThreshold = cf.Duration("slow_threshold", 1*time.Second)

	section, err := cf.Section("query")
	if err != nil {
		panic(err)
	}
	this.Query.loadConfig(section)

	section, err = cf.Section("breaker")
	if err == nil {
		this.Breaker.loadConfig(section)
	}
	this.Servers = make(map[string]*ConfigMysqlInstance)
	for i := 0; i < len(cf.List("servers", nil)); i++ {
		section, err := cf.Section(fmt.Sprintf("servers[%d]", i))
		if err != nil {
			panic(err)
		}

		server := new(ConfigMysqlInstance)
		server.ConnectTimeout = this.ConnectTimeout
		server.loadConfig(section)
		this.Servers[server.Pool] = server
	}

	log.Debug("mysql config: %+v", *this)
}

type ConfigBreaker struct {
	FailureAllowance uint
	RetryTimeout     time.Duration
}

func (this *ConfigBreaker) loadConfig(cf *conf.Conf) {
	this.FailureAllowance = uint(cf.Int("failure_allowance", 5))
	this.RetryTimeout = time.Second * time.Duration(cf.Int("retry_timeout", 10))
}

type ConfigMysqlQuery struct {
	Job   string
	March string
	Pve   string
}

func (this *ConfigMysqlQuery) loadConfig(cf *conf.Conf) {
	this.Job = cf.String("job", "")
	this.March = cf.String("march", "")
	this.Pve = cf.String("pve", "")
	if this.Job == "" &&
		this.March == "" &&
		this.Pve == "" {
		panic("empty mysql query")
	}
}

type ConfigMysqlInstance struct {
	ConnectTimeout time.Duration

	Pool    string
	Host    string
	Port    string
	User    string
	Pass    string
	DbName  string
	Charset string

	dsn string
}

func (this *ConfigMysqlInstance) loadConfig(section *conf.Conf) {
	this.Pool = section.String("pool", "")
	this.Host = section.String("host", "")
	this.Port = section.String("port", "3306")
	this.DbName = section.String("db", "")
	this.User = section.String("username", "")
	this.Pass = section.String("password", "")
	this.Charset = section.String("charset", "utf8")
	if this.Host == "" ||
		this.Port == "" ||
		this.Pool == "" ||
		this.DbName == "" {
		panic("required field missing")
	}

	this.dsn = ""
	if this.User != "" {
		this.dsn = this.User + ":"
		if this.Pass != "" {
			this.dsn += this.Pass
		}
	}
	this.dsn += fmt.Sprintf("@tcp(%s:%s)/%s?", this.Host, this.Port, this.DbName)
	this.dsn += "autocommit=true" // we are not using transaction
	this.dsn += fmt.Sprintf("&timeout=%s", this.ConnectTimeout)
	if this.Charset != "utf8" { // driver default utf-8
		this.dsn += "&charset=" + this.Charset
	}
	this.dsn += "&parseTime=true" // parse db timestamp automatically
}

func (this *ConfigMysqlInstance) String() string {
	return this.DSN()
}

func (this *ConfigMysqlInstance) DSN() string {
	return this.dsn
}
