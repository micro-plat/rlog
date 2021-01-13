package main

import (
	"github.com/micro-plat/hydra"
	"github.com/micro-plat/rlog/services/log"
)

//init 检查应用程序配置文件，并根据配置初始化服务
func init() {

	App.Micro("/log/save", log.NewSaveHandler) //根据配置的日志名称，初始化服务

	hydra.OnReady(func() {

		hydra.Conf.API("7071")
		hydra.Conf.RPC("7011")

		if hydra.IsDebug() {
			hydra.Conf.Vars().Custom("elastic", "logging", `{
				"address": "http://192.168.106.157:9200",
				"write-timeout": 50,
				"cron": "@every 10s",
				"user-name":"elastic",
				"password":"123456"
			}`)
			return
		}
		hydra.Conf.Vars().Custom("elastic", "logging", `{
		"address": "###",
		"write-timeout": 50,
		"cron": "@every 10s"
		}`)
	})

}
