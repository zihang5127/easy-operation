package controller

import (
	"github.com/astaxie/beego"
	"github.com/zihang5127/easy-operation/model"
	"github.com/zihang5127/easy-operation/module/encry"
	"github.com/zihang5127/easy-operation/module/pager"

	"bytes"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego/logs"

	"image/gif"
	"image/jpeg"
	"image/png"
)

// UserController 会员控制器
type UserController struct {
	BaseController
}

// Index 列表
func (c *UserController) Index() {
	c.Prepare()
	c.authenticate()

	if c.User.Role != 0 {
		c.Abort("403")
	}

	c.Layout = ""
	c.TplName = "user/list.html"
	c.Data["UserSelected"] = true

	pageIndex, _ := c.GetInt("page", 1)

	var users []model.User

	pageOptions := pager.PageOptions{
		TableName:           model.NewUser().TableName(),
		EnableFirstLastLink: true,
		CurrentPage:         pageIndex,
		PageSize:            15,
		ParamName:           "page",
		Conditions:          " order by id desc",
	}

	//返回分页信息,
	//第一个:为返回的当前页面数据集合,ResultSet类型
	//第二个:生成的分页链接
	//第三个:返回总记录数
	//第四个:返回总页数
	totalItem, totalCount, rs, pageHtml := pager.GetPagerLinks(&pageOptions, c.Ctx)

	_, err := rs.QueryRows(&users) //把当前页面的数据序列化进一个切片内

	if err != nil {
		logs.Error("", err.Error())
	}

	c.Data["lists"] = users
	c.Data["html"] = pageHtml
	c.Data["totalItem"] = totalItem
	c.Data["totalCount"] = totalCount
}

// 个人中心
func (c *UserController) My() {
	c.Prepare()
	c.authenticate()
	c.Layout = ""
	c.TplName = "user/edit.html"
	c.Data["UserSelected"] = true

	user := c.User

	if c.Ctx.Input.IsPost() {
		password := c.GetString("password")
		status, _ := c.GetInt("status", 0)

		if password != "" {
			pass, _ := encry.PasswordHash(password)
			user.Password = pass
		}
		if user.Role != 0 {
			user.Status = status
		}

		user.Email = c.GetString("email")
		user.Phone = c.GetString("phone")
		user.Avatar = c.GetString("avatar")

		if user.Avatar == "" {
			user.Avatar = "/static/images/headimgurl.jpg"
		}

		var result error

		if user.Id > 0 {
			result = user.Update()
		} else {
			result = user.Add()
		}

		if result != nil {
			c.JsonResult(500, result.Error())
		}

		view, err := c.ExecuteViewPathTemplate("user/list_item.html", *user)
		if err != nil {
			logs.Error("", err.Error())
		}

		data := map[string]interface{}{
			"view": view,
		}
		c.SetUser(*user)
		c.JsonResult(0, "ok", data)

	}

	c.Data["Model"] = user
	c.Data["IsSelf"] = true
}

// Edit 编辑
func (c *UserController) Edit() {
	c.Prepare()
	c.authenticate()

	c.TplName = "user/edit.html"
	userId, _ := c.GetInt(":id")

	user := model.NewUser()

	if userId > 0 {
		user.Id = userId
		if err := user.Find(); err != nil {
			c.ServerError("Data query error:" + err.Error())
		}
	}

	if c.Ctx.Input.IsPost() {
		password := c.GetString("password")
		username := c.GetString("username")
		status, _ := c.GetInt("status", 0)

		if user.Id > 0 {
			if password != "" {
				pass, _ := encry.PasswordHash(password)
				user.Password = pass
			}
			if user.Role != 0 {
				user.Status = status
			}
		} else {
			if username == "" {
				c.JsonResult(500, "Username is require.")
			}
			if password == "" {
				c.JsonResult(500, "Password is require.")
			}
			user.Role = 1
			user.Username = username
			user.Password = password
		}

		user.Email = c.GetString("email")
		user.Phone = c.GetString("phone")
		user.Avatar = c.GetString("avatar")

		if user.Avatar == "" {
			user.Avatar = "/static/images/headimgurl.jpg"
		}

		var result error

		if user.Id > 0 {

			result = user.Update()
		} else {
			result = user.Add()
		}

		if result != nil {
			c.JsonResult(500, result.Error())
		}

		view, err := c.ExecuteViewPathTemplate("user/list_item.html", *user)
		if err != nil {
			logs.Error("", err.Error())
		}

		data := map[string]interface{}{
			"view": view,
		}
		if user.Id == c.User.Id {
			c.SetUser(*user)
		}
		c.JsonResult(0, "ok", data)

	}

	c.Data["Model"] = user
	c.Data["IsSelf"] = false
	if user.Id == c.User.Id {
		c.Data["IsSelf"] = true
	}
}

// 删除用户
func (c *UserController) Delete() {
	c.Prepare()
	c.authenticate()

	if c.User.Role != 0 {
		c.Abort("403")
	}

	userId, err := c.GetInt(":id")

	if err != nil {
		logs.Error("", err.Error())
		c.JsonResult(500, "Parameter error.")
	}

	user := model.NewUser()
	user.Id = userId

	if err := user.Find(); err != nil {
		logs.Error("", err.Error())
		c.JsonResult(500, "Data query error.")
	}

	if user.Role == 0 {
		c.JsonResult(500, "Administrator username cannot be deleted")
	}
	if err := user.Delete(); err != nil {
		logs.Error("", err.Error())
		c.JsonResult(500, "Delete failed")
	}
	c.JsonResult(0, "ok")
}

/*
 * 图片上传
 */
func (c *UserController) Upload() {
	file, moreFile, err := c.GetFile("image-file")
	defer file.Close()

	if err != nil {
		logs.Error("", err.Error())
		c.JsonResult(500, "Read file Fail")
	}

	ext := filepath.Ext(moreFile.Filename)

	if !strings.EqualFold(ext, ".png") && !strings.EqualFold(ext, ".jpg") && !strings.EqualFold(ext, ".gif") && !strings.EqualFold(ext, ".jpeg") {
		c.JsonResult(500, "Unsupported image format")
	}

	x1, _ := strconv.ParseFloat(c.GetString("x"), 10)
	y1, _ := strconv.ParseFloat(c.GetString("y"), 10)
	w1, _ := strconv.ParseFloat(c.GetString("width"), 10)
	h1, _ := strconv.ParseFloat(c.GetString("height"), 10)

	x := int(x1)
	y := int(y1)
	width := int(w1)
	height := int(h1)

	logs.Info("%s %s %s %s", x, x1, y, y1)
	fileName := "avatar_" + strconv.FormatInt(int64(time.Now().Nanosecond()), 16)

	filePath := "static/uploads/" + time.Now().Format("200601") + "/" + fileName + ext

	path := filepath.Dir(filePath)

	_ = os.MkdirAll(path, os.ModePerm)

	err = c.SaveToFile("image-file", filePath)

	if err != nil {
		logs.Error("", err)
		c.JsonResult(500, "Image save failed")
	}

	fileBytes, err := ioutil.ReadFile(filePath)

	if err != nil {
		logs.Error("", err)
		c.JsonResult(500, "Image save failed")
	}

	buf := bytes.NewBuffer(fileBytes)

	m, _, err := image.Decode(buf)

	if err != nil {
		logs.Error("image.Decode => ", err)
		c.JsonResult(500, "Image decoding failed")
	}

	var subImg image.Image

	if rgbImg, ok := m.(*image.YCbCr); ok {
		subImg = rgbImg.SubImage(image.Rect(x, y, x+width, y+height)).(*image.YCbCr) //图片裁剪x0 y0 x1 y1
	} else if rgbImg, ok := m.(*image.RGBA); ok {
		subImg = rgbImg.SubImage(image.Rect(x, y, x+width, y+height)).(*image.YCbCr) //图片裁剪x0 y0 x1 y1
	} else if rgbImg, ok := m.(*image.NRGBA); ok {
		subImg = rgbImg.SubImage(image.Rect(x, y, x+width, y+height)).(*image.YCbCr) //图片裁剪x0 y0 x1 y1
	} else {
		logs.Info("%s",m)
		c.JsonResult(500, "Image decoding failed")
	}

	f, err := os.OpenFile("./"+filePath, os.O_SYNC|os.O_RDWR, 0666)

	if err != nil {
		c.JsonResult(500, "Image save failed")
	}
	defer f.Close()

	if strings.EqualFold(ext, ".jpg") || strings.EqualFold(ext, ".jpeg") {

		err = jpeg.Encode(f, subImg, &jpeg.Options{Quality: 100})
	} else if strings.EqualFold(ext, ".png") {
		err = png.Encode(f, subImg)
	} else if strings.EqualFold(ext, ".gif") {
		err = gif.Encode(f, subImg, &gif.Options{NumColors: 256})
	}
	if err != nil {
		logs.Error("Picture clipping failed => ", err.Error())
		c.JsonResult(500, "Picture clipping failed")
	}

	if err != nil {
		logs.Error("File save failed => ", err.Error())
		c.JsonResult(500, "File save failed")
	}
	url := "/" + filePath

	c.JsonResult(0, "ok", url)
}




// Login 用户登录.
func (c *UserController) Login()  {
	c.Prepare()

	if c.Ctx.Input.IsPost() {
		username := c.GetString("username")
		password := c.GetString("password")

		user,err := model.NewUser().Login(username,password)

		//如果没有数据
		if err == nil {
			c.SetUser(*user)
			c.JsonResult(0,"ok")
			c.StopRun()
		}else{
			logs.Error("%s",err)
			c.JsonResult(500,"Wrong username or password",nil)
		}

		return
	}else{

		c.Layout = ""
		c.TplName = "user/login.html"
	}
}

// Logout 退出登录.
func (c *UserController) Logout(){
	c.SetUser(model.User{})

	c.Redirect(beego.URLFor("UserController.Login"),302)
}
