package main

import (
	"gitee.com/cristiane/micro-mall-pay-consumer/startup"
	"gitee.com/kelvins-io/kelvins"
	"gitee.com/kelvins-io/kelvins/app"
)

const APP_NAME = "micro-mall-pay-consumer"

func main() {
	application := &kelvins.QueueApplication{
		Application: &kelvins.Application{
			LoadConfig: startup.LoadConfig,
			SetupVars:  startup.SetupVars,
			Name:       APP_NAME,
		},
		GetNamedTaskFuncs: startup.GetNamedTaskFuncs,
	}
	app.RunQueueApplication(application)
}
