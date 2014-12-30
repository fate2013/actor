package config

import (
	conf "github.com/funkygao/jsconf"
)

type ConfigBeanstalk struct {
	Breaker ConfigBreaker
}

func (this *ConfigBeanstalk) loadConfig(cf *conf.Conf) {

}
