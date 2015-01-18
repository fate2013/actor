package config

import (
	conf "github.com/funkygao/jsconf"
)

type ConfigWorkerRtm struct {
	MaxProcs int
	Backlog int

	PrimaryHosts []string
	BackupHosts []string
	Timeout int
	ProjectId int
	SecretKey string
}

func (this *ConfigWorkerRtm) loadConfig(cf *conf.Conf) {
	this.MaxProcs = cf.Int("max_procs", 50)
	this.Backlog = cf.Int("backlog", 200)
    primaryHosts := cf.List("primary_hosts", nil)
    for _, host := range primaryHosts {
        this.PrimaryHosts = append(this.PrimaryHosts, host.(string))
    }
    backupHosts := cf.List("backup_hosts", nil)
    for _, host := range backupHosts {
        this.BackupHosts = append(this.BackupHosts, host.(string))
    }
	this.Timeout = cf.Int("timeout", 1000)
	this.ProjectId = cf.Int("project_id", 10001)
	this.SecretKey = cf.String("secret_key", "test_key")
}
