package main

import (
	"core"
	"core/consts/service"
	. "core/libs"
	"core/libs/sessions"
	"core/service"
	"servives/connector/module"
	"servives/public/rpcModules"
)

func main() {
	//初始化Service
	newService := service.NewService(Service.Connector)
	newService.StartRedis()
	newService.StartWebSocket()
	newService.SetSessionCreateHandle(sessionCreate)
	newService.StartIpcClient([]string{Service.Game, Service.Login, Service.Chat})
	newService.StartRpcClient([]string{Service.Game, Service.Login, Service.Chat})

	//模块初始化
	initModule()

	//保持进程
	Run()
}

func initModule() {
	module.StartServerTimer()
}

func sessionCreate(session *sessions.FrontSession) {
	session.AddCloseCallback(nil, "FrontSessionOffline", func() {
		sessionOffline(session)
	})
}

func sessionOffline(session *sessions.FrontSession) {
	method := "Client.Offline"
	args := &rpcModules.ClientOfflineReq{
		ServiceIdentify: core.Service.Identify(),
		UserSessionId:   session.ID(),
	}
	reply := &rpcModules.ClientOfflineReq{}

	//通知登录服务器
	go func() {
		loginService := core.Service.GetRpcClient(Service.Login)
		loginService.Call(method, args, reply, "")
	}()

	//通知聊天服务器
	go func() {
		chatService := core.Service.GetRpcClient(Service.Chat)
		chatService.Call(method, args, reply, "")
	}()
}
