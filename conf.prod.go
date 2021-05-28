// +build prod

package main

import (
	"github.com/micro-plat/hydra"
	"github.com/micro-plat/hydra/conf"
	"github.com/micro-plat/hydra/conf/server/cron"
	"github.com/micro-plat/rlog/services"
)

func init() {
	hydra.Conf.CRON(cron.WithMasterSlave())
	hydra.Conf.API(conf.ByInstall).Metric(conf.ByInstall, conf.ByInstall, "@every 10s")
	hydra.Conf.RPC(conf.ByInstall).Metric(conf.ByInstall, conf.ByInstall, "@every 10s")
	hydra.Conf.Vars().Custom("conf", "clearConf", &services.ClearConf{
		ExpireDay: 15,
	})

	hydra.Conf.Vars().Custom("elastic", "logging", &services.Conf{
		Address:      conf.ByInstall,
		UserName:     "",
		Password:     "",
		WriteTimeout: 50,
		Cron:         10,
	})
}
