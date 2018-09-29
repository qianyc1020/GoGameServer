package main

import (
	"core/consts/service"
	. "core/libs"
	"core/service"
)

func main() {
	//初始化Service
	newService := service.NewService(Service.Log)
	newService.StartRpcServer()

	//模块初始化
	initModule()

	//保持进程
	Run()
}

func initModule() {

}
