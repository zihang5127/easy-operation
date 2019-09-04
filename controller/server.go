package controller

import (
	"bytes"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/zihang5127/easy-operation/model"
	"github.com/zihang5127/easy-operation/module/pager"
	"strconv"
)

// ServerController 服务器控制器
type ServerController struct {
	BaseController
}

// Index 服务器列表
func (c *ServerController) Index() {
	c.Prepare()
	c.authenticate()

	c.Layout = ""
	c.TplName = "server/index.html"

	pageIndex, _ := c.GetInt("page", 1)

	var servers []model.Server

	pageOptions := pager.PageOptions{
		TableName:           model.NewServer().TableName(),
		EnableFirstLastLink: true,
		CurrentPage:         pageIndex,
		PageSize:            15,
		ParamName:           "page",
	}

	totalItem, totalCount, rs, pageHtml := pager.GetPagerLinks(&pageOptions, c.Ctx)

	_, err := rs.QueryRows(&servers) //把当前页面的数据序列化进一个切片内

	if err != nil {
		logs.Error("%s", err.Error())
	}

	c.Data["lists"] = servers
	c.Data["html"] = pageHtml
	c.Data["totalItem"] = totalItem
	c.Data["totalCount"] = totalCount
	c.Data["Server"] = true
}

// Edit 编辑
func (c *ServerController) Edit() {
	c.Prepare()
	c.authenticate()
	c.Layout = ""
	c.TplName = "server/edit.html"

	if c.Ctx.Input.IsPost() {
		id, _ := c.GetInt("id", 0)
		username := c.GetString("username", "")
		serverName := c.GetString("name", "")
		ipAddress := c.GetString("ip", "")
		port, err := c.GetInt("port", 22)
		status, _ := c.GetInt("status", 0)

		if status != 0 && status != 1 {
			status = 0
		}

		if err != nil {
			c.JsonResult(500, "Port Error")
		}
		tag := c.GetString("tag", "")

		key := c.GetString("key", "")

		if serverName == "" {
			c.JsonResult(500, "Server Name is require.")
		}
		if ipAddress == "" {
			c.JsonResult(500, "Server Ip is require.")
		}
		if port <= 0 {
			c.JsonResult(500, "Port is require.")
		}
		if tag != "" {

		}
		if key == "" {
			c.JsonResult(500, "SSH Private Key or Username Password is require.")
		}

		server := model.NewServer()

		if id > 0 {
			server.Id = id
			if err := server.Find(); err != nil {
				c.JsonResult(500, err.Error())
			}
			//如果不是本人创建则返回403
			if server.CreateAt != c.User.Id {
				c.Abort("403")
			}
		}

		server.Username = username
		server.CreateAt = c.User.Id
		server.IpAddress = ipAddress
		server.Name = serverName
		server.Port = port
		server.Tag = tag
		server.PrivateKey = key
		server.Status = status

		if err := server.Save(); err != nil {
			c.JsonResult(500, "Save failed:"+err.Error())
		} else {
			data := make(map[string]interface{}, 5)

			if id <= 0 {

				var buf bytes.Buffer

				viewPath := c.ViewPath

				if c.ViewPath == "" {
					viewPath = beego.BConfig.WebConfig.ViewsPath

				}

				beego.ExecuteViewPathTemplate(&buf, "server/index_list.html", viewPath, server)

				data["view"] = buf.String()
			}

			data["errcode"] = 0
			data["message"] = "ok"

			data["data"] = server

			c.Data["json"] = data
			c.ServeJSON(true)
			c.StopRun()

		}
	}

	id, err := strconv.Atoi(c.Input().Get(":id"))
	if err != nil {
		c.Abort("404")
	}
	server := model.NewServer()
	server.Id = id

	if err := server.Find(); err != nil {
		c.Abort("404")
	}
	//如果不是本人创建则返回403
	if server.CreateAt != c.User.Id {
		c.Abort("403")
	}
	if c.Ctx.Input.IsAjax() {

		c.JsonResult(0, "ok", *server)
	}
	c.Data["Model"] = server
	c.Data["Server"] = true

}

// Delete 删除一个Server
func (c *ServerController) Delete() {
	c.authenticate()
	id, _ := c.GetInt("id", 0)

	if id <= 0 {
		c.JsonResult(500, "Server ID is require.")
	}
	server := model.NewServer()

	server.Id = id

	if err := server.Find(); err != nil {
		c.JsonResult(500, err.Error())
	}
	if server.CreateAt != c.User.Id {
		c.JsonResult(403, "Permission denied")
	}
	if err := server.Delete(); err != nil {
		c.JsonResult(500, err.Error())
	}

	model.NewPsRelation().DeleteByWhere(" AND server_id = ?", id)

	c.JsonResult(0, "ok")
}
