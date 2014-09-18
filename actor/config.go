package actor

import (
	"errors"
	"fmt"
	conf "github.com/funkygao/jsconf"
	log "github.com/funkygao/log4go"
	"time"
)

type ActorConfig struct {
	RestListenAddr       string
	ProfListenAddr       string
	CallbackUrl          string
	CallbackTimeout      int
	ScheduleInterval     int
	ConsoleStatsInterval int
	MysqlConfig          *ConfigMysql
}

func (this *ActorConfig) LoadConfig(cf *conf.Conf) (err error) {
	this.RestListenAddr = cf.String("rest_listen_addr", "")
	this.ProfListenAddr = cf.String("prof_listen_addr", "")
	this.CallbackUrl = cf.String("callback_url", "")
	if this.CallbackUrl == "" {
		err = errors.New("empty callback_url")
		return
	}
	this.CallbackTimeout = cf.Int("callback_timeout", 4)
	this.ScheduleInterval = cf.Int("sched_interval", 1)
	this.ConsoleStatsInterval = cf.Int("stats_interval", 60*10)

	this.MysqlConfig = new(ConfigMysql)
	var section *conf.Conf
	section, err = cf.Section("mysql")
	if err != nil {
		return
	}
	this.MysqlConfig.loadConfig(section)

	log.Debug("actor config %+v", *this)
	return
}

type ConfigMysql struct {
	ConnectTimeout time.Duration
	IoTimeout      time.Duration
	SlowThreshold  time.Duration
	Breaker        ConfigBreaker
	Servers        map[string]*ConfigMysqlInstance // key is pool
}

func (this *ConfigMysql) Pools() (pools []string) {
	poolsMap := make(map[string]bool)
	for _, server := range this.Servers {
		poolsMap[server.Pool] = true
	}
	for poolName, _ := range poolsMap {
		pools = append(pools, poolName)
	}
	return
}

func (this *ConfigMysql) loadConfig(cf *conf.Conf) {
	this.ConnectTimeout = time.Duration(cf.Int("connect_timeout", 4)) * time.Second
	this.IoTimeout = time.Duration(cf.Int("io_timeout", 30)) * time.Second
	this.SlowThreshold = time.Duration(cf.Int("slow_threshold", 5)) * time.Second
	section, err := cf.Section("breaker")
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

type ConfigMysqlInstance struct {
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
	if this.Charset != "" {
		this.dsn += "charset=" + this.Charset
	}

	log.Debug("mysql instance: %s", this.dsn)
}

func (this *ConfigMysqlInstance) DSN() string {
	return this.dsn
}
