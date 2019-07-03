package main

import (
	"github.com/micro-plat/hydra/hydra"
)

type logSaver struct {
	*hydra.MicroApp
}

func main() {
	app := &logSaver{
		hydra.NewApp(
			hydra.WithPlatName("logging"),
			hydra.WithSystemName("logsaver"),
			hydra.WithServerTypes("rpc-api"),
			hydra.WithDebug()),
	}

	app.init()
	app.install()

	app.Start()
}
