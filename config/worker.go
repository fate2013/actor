package config

import (
	conf "github.com/funkygao/jsconf"
)

type ConfigWorker struct {
	Php ConfigWorkerPhp
	Pnb ConfigWorkerPnb
}

func (this *ConfigWorker) loadConfig(cf *conf.Conf) {
	section, err := cf.Section("php")
	if err != nil {
		panic(err)
	}
	this.Php.loadConfig(section)

	section, err = cf.Section("pnb")
	if err != nil {
		panic(err)
	}
	this.Pnb.loadConfig(section)
}
