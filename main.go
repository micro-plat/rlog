package main

import (
	"github.com/micro-plat/hydra"
	"github.com/micro-plat/hydra/hydra/servers/cron"
	"github.com/micro-plat/hydra/hydra/servers/http"
	"github.com/micro-plat/hydra/hydra/servers/rpc"
)

var App = hydra.NewApp(
	hydra.WithPlatName("logging"),
	hydra.WithSystemName("logsaver"),
	hydra.WithServerTypes(http.API, rpc.RPC, cron.CRON))

func main() {
	App.Start()
}
