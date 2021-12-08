package main

import (
	"github.com/micro-plat/hydra/conf/app"
	"github.com/micro-plat/rlog/services"
)

//init 检查应用程序配置文件，并根据配置初始化服务
func init() {
	App.Micro("/log/save", services.SaveHandle)
	App.Micro("/buried/save", services.BuriedHandle)

	App.CRON("/log/clear", services.ClearHandle, "@every 10h")

	App.OnStarting(func(appconf app.IAPPConf) error {
		return nil
	})
}
