package router

import (
	"github.com/astaxie/beego"
	"github.com/zihang5127/easy-operation/controller"
)

func init() {
	beego.Router("/*", &controller.IndexController{}, "*:Index")
	beego.Router("/login", &controller.UserController{}, "*:Login")
	beego.Router("/logout", &controller.UserController{}, "*:Logout")
	beego.Router("/user/list", &controller.UserController{}, "*:Index")
	beego.Router("/user/edit:id", &controller.UserController{}, "*:Edit")
	beego.Router("/user/delete:id", &controller.UserController{}, "*:Delete")
	beego.Router("/my", &controller.UserController{}, "*:My")
	beego.Router("/user/upload", &controller.UserController{}, "post:Upload")

	beego.Router("/project/list", &controller.IndexController{}, "*:Index")
	beego.Router("/project/edit:id", &controller.IndexController{}, "*:Edit")
	beego.Router("/project/delete:id", &controller.IndexController{}, "*:Delete")
	beego.Router("/project/build", &controller.IndexController{}, "*:Build")
	beego.Router("/project/server/list:id", &controller.IndexController{}, "*:ServerList")
	beego.Router("/project/server/delete:id", &controller.IndexController{}, "*:DeleteServer")
	beego.Router("/project/server/add", &controller.IndexController{}, "*:AddServer")

	beego.Router("/server/list", &controller.ServerController{}, "*:Index")
	beego.Router("/server/edit:id", &controller.ServerController{}, "*:Edit")
	beego.Router("/server/delete:id", &controller.ServerController{}, "*:Delete")

	beego.Router("/ws", &controller.WebSocketController{}, "*:Ws")
}
