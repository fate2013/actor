package config

import (
	conf "github.com/funkygao/jsconf"
)

type ConfigPubnub struct {
	PublishKey   string
	SubscribeKey string
	SecretKey    string
	Cipher       string
	UseSsl       bool
}

func (this *ConfigPubnub) loadConfig(cf *conf.Conf) {
	this.PublishKey = cf.String("publish_key", "")
	this.SubscribeKey = cf.String("subscribe_key", "")
	this.SecretKey = cf.String("secret_key", "")
	this.Cipher = cf.String("cipher", "")
	this.UseSsl = cf.Bool("use_ssl", false)
}
