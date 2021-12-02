//go:build !prod
// +build !prod

package main

import (
	"github.com/micro-plat/hydra"
	"github.com/micro-plat/hydra/conf/server/cron"
	"github.com/micro-plat/rlog/services"
)

func init() {
	hydra.Conf.CRON(cron.WithMasterSlave())
	hydra.Conf.API("7071").Metric("http://192.168.0.185:8086", "convoy", "@every 10s")
	hydra.Conf.RPC("7011").Metric("http://192.168.0.185:8086", "convoy", "@every 10s")
	hydra.Conf.Vars().Custom("conf", "clearConf", &services.ClearConf{
		ExpireDay: 15,
	})

	hydra.Conf.Vars().Custom("elastic", "logging", &services.Conf{
		Address:      `http://192.168.106.177:9200,http://192.168.0.126:9200,http://192.168.0.125:9200`,
		UserName:     "",
		Password:     "",
		WriteTimeout: 50,
		Cron:         10,
	})
	return
}
