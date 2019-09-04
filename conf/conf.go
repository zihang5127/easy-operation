package conf

import "github.com/astaxie/beego"

func QueueSize() int {
	return beego.AppConfig.DefaultInt("queue_size", 100)
}
