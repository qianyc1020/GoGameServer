package main

import (
	"core/consts/service"
	. "core/libs"
	"core/messages"
	"core/protos/gameProto"
	"core/service"
	"servives/login/module"
)

func main() {
	//初始化Service
	newService := service.NewService(Service.Login)
	newService.StartIpcServer()
	newService.StartRpcServer(&module.LoginRpcServer{})
	newService.StartRpcClient([]string{Service.Log})
	newService.StartRedis()
	newService.StartMysql()

	//消息初始化
	initMessage()

	//模块初始化
	initModule()

	//保持进程
	Run()
}

func initMessage() {
	messages.RegisterIpcServerHandle(gameProto.ID_user_login_c2s, module.Login)
}

func initModule() {

}