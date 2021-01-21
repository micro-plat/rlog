package main

import (
	"github.com/micro-plat/hydra"
	"github.com/micro-plat/hydra/conf"
	"github.com/micro-plat/rlog/services"
)

//init 检查应用程序配置文件，并根据配置初始化服务
func init() {

	App.Micro("/log/save", services.SaveHandle) //根据配置的日志名称，初始化服务

	hydra.OnReady(func() {

		hydra.Conf.API("7071")
		hydra.Conf.RPC("7011")

		if hydra.IsDebug() {
			hydra.Conf.Vars().Custom("elastic", "logging", &services.Conf{
				Address:      "http://192.168.106.177:9200",
				UserName:     "elastic",
				Password:     "123456",
				WriteTimeout: 50,
				Cron:         10,
			})
			return
		}

		hydra.Conf.Vars().Custom("elastic", "logging", &services.Conf{
			Address:      conf.ByInstall,
			UserName:     "",
			Password:     "",
			WriteTimeout: 50,
			Cron:         10,
		})
	})
}
