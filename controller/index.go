package controller

import (
	"github.com/astaxie/beego/logs"
	_ "github.com/astaxie/beego/logs"
	"github.com/zihang5127/easy-operation/model"
	"github.com/zihang5127/easy-operation/module/pager"
	"github.com/zihang5127/easy-operation/module/tasks"
	"strings"

	"strconv"
)

// indexController 首页Project控制器
type IndexController struct {
	BaseController
}

// Index 首页
func (c *IndexController) Index() {

	c.Prepare()
	c.authenticate()
	c.Layout = ""
	c.TplName = "index/index.html"

	pageIndex, _ := c.GetInt("page", 1)

	var projectes []model.Project

	pageOptions := pager.PageOptions{
		TableName:           model.NewProject().TableName(),
		EnableFirstLastLink: true,
		CurrentPage:         pageIndex,
		ParamName:           "page",
		PageSize:            15,
		Conditions:          " AND create_at = 1 order by id desc",
	}

	//返回分页信息,
	//第一个:为返回的当前页面数据集合,ResultSet类型
	//第二个:生成的分页链接
	//第三个:返回总记录数
	//第四个:返回总页数
	totalItem, totalCount, rs, pageHtml := pager.GetPagerLinks(&pageOptions, c.Ctx)

	_, err := rs.QueryRows(&projectes) //把当前页面的数据序列化进一个切片内

	if err != nil {
		logs.Error("%s", err.Error())
	}

	c.Data["lists"] = projectes
	c.Data["html"] = pageHtml
	c.Data["totalItem"] = totalItem
	c.Data["totalCount"] = totalCount
}

// Edit 编辑
func (c *IndexController) Edit() {
	c.Prepare()
	c.authenticate()
	c.Layout = ""
	c.TplName = "index/edit.html"

	if c.Ctx.Input.IsPost() {
		id, _ := c.GetInt("id", 0)
		name := strings.TrimSpace(c.GetString("repo_name", ""))
		branch := strings.TrimSpace(c.GetString("branch_name", ""))
		tag := strings.TrimSpace(c.GetString("tag", ""))
		shell := strings.TrimSpace(c.GetString("shell", ""))
		status, _ := c.GetInt("status", 0)
		repositoryType := c.GetString("repository_type", "")
		if name == "" {
			c.JsonResult(500, "Repository Name is require.")
		}
		if branch == "" {
			branch = "master"
		}

		if shell == "" {
			c.JsonResult(500, "Callback Shell Script is require.")
		}

		project := model.NewProject()

		if id > 0 {
			project.Id = id
			if err := project.Find(); err != nil {
				c.JsonResult(500, err.Error())
			}
			if project.CreateAt != c.User.Id {
				c.JsonResult(403, "Permission denied")
			}
		}

		project.RepositoryName = name
		project.BranchName = branch
		project.Tag = tag
		project.Shell = shell
		project.Status = status
		project.CreateAt = c.User.Id
		project.RepositoryType = repositoryType

		if err := project.Save(); err != nil {
			c.JsonResult(500, err.Error())
		}
		data := make(map[string]interface{}, 5)

		if id <= 0 {
			view, _ := c.ExecuteViewPathTemplate("index/index_list.html", project)
			data["view"] = view
		}

		data["errcode"] = 0
		data["message"] = "ok"

		data["data"] = project

		c.Data["json"] = data
		c.ServeJSON(true)
		c.StopRun()

	}

	id, _ := strconv.Atoi(c.Input().Get(":id"))

	if id <= 0 {
		c.Abort("404")
	}

	project := model.NewProject()
	project.Id = id

	if err := project.Find(); err != nil {
		c.TplName = "errors/500.html"
		c.Data["Message"] = err.Error()
	} else {
		c.Data["Model"] = project
	}
}

// Delete 删除
func (c *IndexController) Delete() {
	c.authenticate()
	id, _ := c.GetInt("id", 0)
	if id <= 0 {
		c.JsonResult(500, "Server ID is require.")
	}

	project := model.NewProject()
	project.Id = id

	if err := project.Find(); err != nil {
		c.JsonResult(500, "Git Project does not exist")
	}
	if project.CreateAt != c.User.Id {
		c.JsonResult(403, "Permission denied")
	}

	if err := project.Delete(); err != nil {
		c.JsonResult(500, "failed to delete")
	}

	model.NewPsRelation().DeleteByWhere(" AND project_id = ?", id)

	c.JsonResult(0, "ok")
}

/**
 * 执行构建
 */
func (c *IndexController) Build() {
	c.Prepare()
	c.authenticate()
	c.TplName = "index/log.html"

	id, err := strconv.Atoi(c.Input().Get(":id"))

	if err != nil || id <= 0 {
		c.ServerError("Project does not exist.")
	}

	project := model.NewProject()
	project.Id = id

	if err := project.Find(); err != nil {
		c.TplName = "errors/500.html"
		c.Data["Message"] = err.Error()
	} else {
		if project.Status != 0 {
			logs.Error(500, "Project disabled.")
			c.ServerError("Project disabled.")
			c.StopRun()
		}
		c.Data["Model"] = project
	}

	psRelationDetailed, err := model.FindRelationDetailedByWhere("AND ps.project_id = ?", project.Id)

	if err != nil {
		logs.Error(500, err.Error())
		c.ServerError("Build error")
	}
	if len(psRelationDetailed) <= 0 {
		logs.Error(500, "Servers empty or disabled.")
		c.ServerError("Servers empty or disabled.")
	}

	for _, job := range psRelationDetailed {
		tasks.Add(tasks.Task{ServerId: job.ServerId, ProjectId: job.ProjectId})
	}
}

func (c *IndexController) ServerList() {
	c.Prepare()
	c.authenticate()
	c.TplName = "index/server_list.html"

	id, err := strconv.Atoi(c.Input().Get(":id"))
	if err != nil || id <= 0 {
		c.ServerError("Project does not exist.")
	}
	project := model.NewProject()
	project.Id = id

	if err := project.Find(); err != nil {
		c.ServerError("Project does not exist.")
	}
	if project.CreateAt != c.User.Id {
		c.Forbidden("")
	}

	c.Data["Model"] = project

	res, err := model.NewPsRelation().QueryByProjetId(id, c.User.Id)

	c.Data["lists"] = res
}

// AddServer 检索服务器并添加到数据库
func (c *IndexController) AddServer() {
	c.Prepare()
	c.authenticate()

	if c.Ctx.Input.IsPost() {
		id, err := c.GetInt("project_id")

		if err != nil {
			c.JsonResult(500, "Parameter error: project_id is require.")
		}
		serverParams := c.GetStrings("server_id")
		if len(serverParams) <= 0 {
			c.JsonResult(500, "Server Id is require.")
		}

		project := model.NewProject()
		project.Id = id

		if err := project.Find(); err != nil {
			c.JsonResult(404, "Project does not exist.")
		}
		if project.CreateAt != c.User.Id {
			c.JsonResult(403, "Permission denied")
		}

		serverIds := make([]int, len(serverParams))
		index := 0
		for _, id := range serverParams {
			if id, err := strconv.Atoi(id); err == nil {
				serverIds[index] = id
				index++
			}

		}
		servers, err := model.NewServer().QueryServerByServerId(serverIds, c.User.Id)

		if err != nil {
			c.JsonResult(500, "An error occurred while querying data")
		}

		if len(servers) <= 0 {
			c.JsonResult(500, "Invalid server")
		}

		relations := make([]map[string]interface{}, len(servers))

		index = 0
		for _, server := range servers {
			ps := model.NewPsRelation()

			ps.ProjectId = id
			ps.ServerId = server.Id
			ps.UserId = c.User.Id

			if err := ps.Save(); err == nil {
				temp := map[string]interface{}{
					"server_id":   server.Id,
					"name":        server.Name,
					"ip_address":  server.IpAddress,
					"port":        server.Port,
					"add_time":    ps.CreateTime,
					"relation_id": ps.ProjectId,
					"status":      server.Status,
				}
				relations[index] = temp
				index++
			}
		}

		c.JsonResult(0, "ok", relations)
	}

	keyword := c.GetString("keyword", "")

	if keyword == "" {
		c.JsonResult(500, "Keyword is require.")
	}

	id, _ := c.GetInt("id")

	ps := model.NewPsRelation()

	var serverIds []int
	if pss, err := ps.QueryByProjetId(id, c.User.Id); err == nil && len(pss) > 0 {
		serverIds = make([]int, len(pss))
		i := 0
		for _, item := range pss {
			serverIds[i] = item.ServerId
			i++
		}
	}

	serverList, err := model.NewServer().Search(keyword, c.User.Id, serverIds...)

	if err != nil {
		c.JsonResult(500, "Query Result Error")
	}

	c.JsonResult(0, "ok", serverList)

	c.StopRun()
}

// DeleteServer 删除一个服务器
func (c *IndexController) DeleteServer() {
	c.Prepare()
	ps_id, err := strconv.Atoi(c.Input().Get(":id"))

	if err != nil {
		c.JsonResult(500, "Parameter error: project_id is require.")
	}

	ps := model.NewPsRelation()

	if err := ps.Find(ps_id); err != nil {
		logs.Info("Delete server failed: %s", err)

		c.JsonResult(404, "Server does not exist.")
	}

	server := model.NewServer()
	server.Id = ps.ServerId

	if err := server.Find(); err != nil || server.CreateAt != c.User.Id {
		c.JsonResult(403, "Permission denied")
	}
	project := model.NewProject()
	project.Id = ps.ProjectId

	if err := project.Find(); err != nil || project.CreateAt != c.User.Id {
		c.JsonResult(403, "Permission denied")
	}

	if err := ps.Delete(); err != nil {
		c.JsonResult(500, "Delete failed")
	}

	c.JsonResult(0, "ok")
}

func (c *IndexController) Log() {
	c.Prepare()
	c.authenticate()
	c.Layout = ""
	c.TplName = "index/log.html"
}
