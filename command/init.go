package command

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/zihang5127/easy-operation/model"
	"net/url"

	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"time"
)

// 初始化数据库
func InitDatabase() {
	host := beego.AppConfig.String("db_host")
	database := beego.AppConfig.String("db_database")
	username := beego.AppConfig.String("db_username")
	password := beego.AppConfig.String("db_password")
	timezone := beego.AppConfig.String("timezone")
	port := beego.AppConfig.String("db_port")
	dataSource := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=true&loc=%s", username, password, host, port, database, url.QueryEscape(timezone))

	_ = orm.RegisterDataBase("default", "mysql", dataSource)

	orm.DefaultTimeLoc, _ = time.LoadLocation(timezone)
}

// 初始化日志
func InitLogger() {

	_ = logs.SetLogger("console")
	_ = logs.SetLogger("file", `{"filename":"logs/log.log"}`)
	logs.EnableFuncCallDepth(true)
	logs.Async()
}

// 初始化模块
func InitModel() {
	orm.RegisterModel(new(model.User))
	orm.RegisterModel(new(model.Project))
	orm.RegisterModel(new(model.Server))
	orm.RegisterModel(new(model.PsRelation))
}

// 初始化命令行
func InitCommand() {
	orm.RunCommand()
	Install()
	Version()
}

// 启动
func Run() {
	beego.Run()
}
