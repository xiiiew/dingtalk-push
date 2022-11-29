package main

import (
	"dingtalk-push/bootstrap"
	"dingtalk-push/conf"
	"dingtalk-push/web/router"
	"fmt"
	"github.com/kataras/iris/v12"
)

func main() {
	app := bootstrap.New("dingtalk-push", "xiiiew")
	app.Bootstrapper(router.RegisterRouter)
	app.Listen(fmt.Sprintf(":%d", conf.ConfigYamlInstance.HTTPConfig.ServerListenPort), iris.WithoutBodyConsumptionOnUnmarshal)
}
