package config

import (
	conf "github.com/funkygao/jsconf"
)

type ConfigPoller struct {
	Mysql     ConfigMysql
	Beanstalk ConfigBeanstalk
}

func (this *ConfigPoller) loadConfig(cf *conf.Conf) {
	section, err := cf.Section("mysql")
	if err != nil {
		panic(err)
	}
	this.Mysql.loadConfig(section)

	section, err = cf.Section("beanstalk")
	if err != nil {
		panic(err)
	}
	this.Beanstalk.loadConfig(section)
}
