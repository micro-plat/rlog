// +build prod

package main

//bindConf 绑定启动配置， 启动时检查注册中心配置是否存在，不存在则引导用户输入配置参数并自动创建到注册中心
func (app *logSaver) install() {
	app.IsDebug = false
	//app.Conf.SetInput("#name", "系统名称", "需要记录日志的系统名称")
	app.Conf.SetInput("#address", "elastic search服务器地址", "http://host:port")

	app.Conf.API.SetMainConf(`{"address":":7010"}`)
	app.Conf.RPC.SetMainConf(`{"address":":7011"}`)

	app.Conf.Plat.SetVarConf("elastic", "logging", `{
		"address": "#address",
		"write-timeout": 50,
		"cron": "@every 10s"
	}`)
}
