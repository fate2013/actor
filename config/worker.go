package config

import (
	conf "github.com/funkygao/jsconf"
)

type ConfigWorker struct {
	Php ConfigWorkerPhp
	Pnb ConfigWorkerPnb
}

func (this *ConfigWorker) loadConfig(cf *conf.Conf) {

}
