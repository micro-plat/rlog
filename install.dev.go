// +build !prod

package main

//bindConf 绑定启动配置， 启动时检查注册中心配置是否存在，不存在则引导用户输入配置参数并自动创建到注册中心
func (app *logSaver) install() {
	//app.Conf.API.SetMainConf(`{"address":":7801"}`)
	app.Conf.RPC.SetMainConf(`{"address":":7802"}`)
	app.Conf.Plat.SetVarConf("elastic", "logging", `{
		"address": "http://192.168.106.157:9200",
		"write-timeout": 50,
		"cron": "@every 10s",
		"user-name":"elastic",
		"password":"123456"
	}`)
}
