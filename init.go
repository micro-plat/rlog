package main

import (
	"github.com/micro-plat/hydra/component"
	"github.com/micro-plat/logsaver/modules/logging"
	"github.com/micro-plat/logsaver/services/log"
)

//init 检查应用程序配置文件，并根据配置初始化服务
func (app *logSaver) init() {

	app.Micro("/log/save", log.NewSaveHandler) //根据配置的日志名称，初始化服务

	app.Closing(func(c component.IContainer) error {
		return logging.Close()
	})
}
