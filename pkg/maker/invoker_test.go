package maker

import (
	"encoding/json"
	"github.com/team-ide/go-tool/util"
	"go.uber.org/zap"
	"sync"
	"teamide/pkg/maker/invokers"
	"teamide/pkg/maker/modelers"
	"testing"
)

func TestInvoker(t *testing.T) {
	app, err := LoadDemoApp()
	if err != nil {
		util.Logger.Error("load demo app error", zap.Error(err))
		return
	}

	invoker := invokers.NewInvoker(app)

	invokeData := invokers.NewInvokeData(app)

	err = invokeData.AddArg("userId", 1, modelers.ValueTypeInt64)
	if err != nil {
		util.Logger.Error("invoke data add arg error", zap.Error(err))
		return
	}

	serviceName := "user/get"
	res, err := invoker.InvokeServiceByName(serviceName, invokeData)
	if err != nil {
		util.Logger.Error("service invoke error", zap.Any("serviceName", serviceName), zap.Error(err))
		return
	}
	bs, err := json.Marshal(res)
	if err != nil {
		util.Logger.Error("res to json error", zap.Error(err))
		return
	}
	println("service [" + serviceName + "] run success")
	println(string(bs))
}

func TestInvokerZk(t *testing.T) {
	app, err := LoadDemoApp()
	if err != nil {
		util.Logger.Error("load demo app error", zap.Error(err))
		return
	}

	invoker := invokers.NewInvoker(app)

	invokeData := invokers.NewInvokeData(app)

	serviceName := "task/zk"
	res, err := invoker.InvokeServiceByName(serviceName, invokeData)
	if err != nil {
		util.Logger.Error("service invoke error", zap.Any("serviceName", serviceName), zap.Error(err))
		return
	}
	bs, err := json.Marshal(res)
	if err != nil {
		util.Logger.Error("res to json error", zap.Error(err))
		return
	}
	println("service [" + serviceName + "] run success")
	println(string(bs))

	wait := sync.WaitGroup{}
	wait.Add(1)
	wait.Wait()
}
