package pager

import (
	// "fmt"
	"github.com/astaxie/beego/context"
	"github.com/astaxie/beego/orm"
	html "html/template"
	con "strconv"
	"strings"
	"time"
)

type PageOptions struct {
	TableName           string //表名  -----------------[必填]
	Conditions          string //条件
	CurrentPage         int    //当前页 ,默认1 每次分页,必须在前台设置新的页数,不设置始终默认1.在控制器中使用方式:cp, _ := this.GetInt("pno")   po.CurrentPage = int(cp)
	PageSize            int    //页面大小,默认20
	LinkItemCount       int    //生成A标签的个数 默认10个
	Href                string //A标签的链接地址  ---------[不需要设置]
	ParamName           string //参数名称  默认是pno
	FirstPageText       string //首页文字  默认"首页"
	LastPageText        string //尾页文字  默认"尾页"
	PrePageText         string //上一页文字 默认"上一页"
	NextPageText        string //下一页文字 默认"下一页"
	EnableFirstLastLink bool   //是否启用首尾连接 默认false 建议开启
	EnablePreNexLink    bool   //是否启用上一页,下一页连接 默认false 建议开启
}

/**
 * 分页函数，适用任何表
 * 返回 总记录条数,总页数,以及当前请求的数据RawSeter,调用中需要"rs.QueryRows(&tblog)"就行了  --tblog是一个Tb_log对象
 * 参数：表名，当前页数，页面大小，条件（查询条件,格式为 " and name='zhifeiya' and age=12 "）
 */
func GetPagesInfo(tableName string, currentpage int, pagesize int, conditions string) (int, int, orm.RawSeter) {
	if currentpage <= 1 {
		currentpage = 1
	}
	if pagesize == 0 {
		pagesize = 20
	}
	var rs orm.RawSeter
	o := orm.NewOrm()
	var totalItem, totalPages = 0, 0                                                                    //总条数,总页数
	_ = o.Raw("SELECT count(*) FROM " + tableName + "  where 1 > 0 " + conditions).QueryRow(&totalItem) //获取总条数
	if totalItem <= pagesize {
		totalPages = 1
	} else if totalItem > pagesize {
		temp := totalItem / pagesize
		if (totalItem % pagesize) != 0 {
			temp = temp + 1
		}
		totalPages = temp
	}
	rs = o.Raw("select *  from  " + tableName + "  where 1 > 0 " + conditions + " LIMIT " + con.Itoa((currentpage-1)*pagesize) + "," + con.Itoa(pagesize))
	return totalItem, totalPages, rs
}

/**
 * 返回总记录条数,总页数,当前页面数据,分页html
 * 根据分页选项,生成分页连接 下面是一个实例:
     func (this *MainController) Test() {
        var po util.PageOptions
        po.EnablePreNexLink = true
        po.EnableFirstLastLink = true
        po.LinkItemCount = 7
        po.TableName = "help_topic"
        cp, _ := this.GetInt("pno")
        po.CurrentPage = int(cp)
        _,_,_ pager := util.GetPagerLinks(&po, this.Ctx)
        this.Data["Email"] = html.HTML(pager)
        this.TplName = "test.html"
    }
*/
func GetPagerLinks(po *PageOptions, ctx *context.Context) (int, int, orm.RawSeter, html.HTML) {
	str := ""
	totalItem, totalpages, rs := GetPagesInfo(po.TableName, po.CurrentPage, po.PageSize, po.Conditions)
	po = setDefault(po, totalpages)
	DealUri(po, ctx)
	if totalpages <= po.LinkItemCount {
		str = fun1(po, totalpages) //显示完全  12345678910
	} else if totalpages > po.LinkItemCount {
		if po.CurrentPage < po.LinkItemCount {
			str = fun2(po, totalpages) //123456789...200
		} else {
			if po.CurrentPage+po.LinkItemCount < totalpages {
				str = fun3(po, totalpages)
			} else {
				str = fun4(po, totalpages)
			}
		}
	}
	return totalItem, totalpages, rs, html.HTML(str)
}

/**
 * 处理url,目的是保存参数
 */
func DealUri(po *PageOptions, ctx *context.Context) {
	uri := ctx.Request.RequestURI
	var rs string
	if strings.Contains(uri, "?") {
		arr := strings.Split(uri, "?")
		rs = arr[0] + "?" + po.ParamName + "time=" + con.Itoa(time.Now().Second())
		arr2 := strings.Split(arr[1], "&")
		for _, v := range arr2 {
			if !strings.Contains(v, po.ParamName) {
				rs += "&" + v
			}
		}
	} else {
		rs = uri + "?" + po.ParamName + "time=" + con.Itoa(time.Now().Second())
	}
	po.Href = rs
}

/**
 * 1...197 198 199 200
 */
func fun4(po *PageOptions, totalPages int) string {
	rs := ""
	rs += getHeader(po, totalPages)
	rs += "<li><a href='" + po.Href + "&" + po.ParamName + "=" + con.Itoa(1) + "'>" + con.Itoa(1) + "</a><li>"
	rs += "<li><a href=''>...</a></li>"
	for i := totalPages - po.LinkItemCount; i <= totalPages; i++ {
		if po.CurrentPage != i {
			rs += "<li><a href='" + po.Href + "&" + po.ParamName + "=" + con.Itoa(i) + "'>" + con.Itoa(i) + "</a></li>"
		} else {
			rs += "<li class=\"active\"><a href=\"###\">" + con.Itoa(i) + "  <span class=\"sr-only\">(current)</span></a></li>"
		}
	}
	rs += getFooter(po, totalPages)
	return rs

}

/**
 * 1...6 7 8 9 10 11 12  13  14 15... 200
 */
func fun3(po *PageOptions, totalpages int) string {
	var rs = ""
	rs += getHeader(po, totalpages)
	rs += "<li><a href='" + po.Href + "&" + po.ParamName + "=" + con.Itoa(1) + "'>" + con.Itoa(1) + "</a></li>"
	rs += "<a href=''>...</a>"
	for i := po.CurrentPage - po.LinkItemCount/2 + 1; i <= po.CurrentPage+po.LinkItemCount/2-1; i++ {
		if po.CurrentPage != i {
			rs += "<a href='" + po.Href + "&" + po.ParamName + "=" + con.Itoa(i) + "'>" + con.Itoa(i) + "</a>"
		} else {
			rs += "<li class=\"active\"><a href=\"###\">" + con.Itoa(i) + "  <span class=\"sr-only\">(current)</span></a></li>"
		}
	}
	rs += "<a href=''>...</a>"
	rs += "<li><a href='" + po.Href + "&" + po.ParamName + "=" + con.Itoa(totalpages) + "'>" + con.Itoa(totalpages) + "</a></li>"
	rs += getFooter(po, totalpages)
	return rs

}

/**
 * totalpages > po.LinkItemCount   po.CurrentPage < po.LinkItemCount
 * 123456789...200
 */
func fun2(po *PageOptions, totalpages int) string {
	var rs = ""
	rs += getHeader(po, totalpages)
	for i := 1; i <= po.LinkItemCount+1; i++ {
		if i == po.LinkItemCount {
			rs += "<li><a href=\"" + po.Href + "&" + po.ParamName + "=" + con.Itoa(i) + "\">...</a></li>"
		} else if i == po.LinkItemCount+1 {
			rs += "<li><a href=\"" + po.Href + "&" + po.ParamName + "=" + con.Itoa(totalpages) + "\">" + con.Itoa(totalpages) + "</a></li>"
		} else {
			if po.CurrentPage != i {
				rs += "<li><a href='" + po.Href + "&" + po.ParamName + "=" + con.Itoa(i) + "'>" + con.Itoa(i) + "</a></li>"
			} else {
				rs += "<li class=\"active\"><a href=\"###\">" + con.Itoa(i) + "  <span class=\"sr-only\">(current)</span></a></li>"
			}
		}
	}
	rs += getFooter(po, totalpages)
	return rs
}

/**
 * totalpages <= po.LinkItemCount
 * 显示完全  12345678910
 */
func fun1(po *PageOptions, totalpages int) string {

	var rs = ""
	rs += getHeader(po, totalpages)
	for i := 1; i <= totalpages; i++ {
		if po.CurrentPage != i {
			rs += "<li><a href='" + po.Href + "&" + po.ParamName + "=" + con.Itoa(i) + "'>" + con.Itoa(i) + "</a></li>"
		} else {
			rs += "<li class=\"active\"><a href=\"###\">" + con.Itoa(i) + "  <span class=\"sr-only\">(current)</span></a></li>"
		}
	}
	rs += getFooter(po, totalpages)
	return rs
}

/**
 * 头部
 */
func getHeader(po *PageOptions, totalpages int) string {
	var rs = "<ul class=\"pagination\">"
	if po.EnableFirstLastLink { //当首页,尾页都设定的时候,就显示

		rs += "<li" + judgeDisable(po, totalpages, 0) + " ><a href='" + po.Href + "&" + po.ParamName + "=" + con.Itoa(1) + "'>" + po.FirstPageText + "</a></li>"
	}
	if po.EnablePreNexLink { // disabled=\"disabled\"
		var a = po.CurrentPage - 1
		if po.CurrentPage == 1 {
			a = 1
		}
		rs += "<li" + judgeDisable(po, totalpages, 0) + " ><a href='" + po.Href + "&" + po.ParamName + "=" + con.Itoa(a) + "'>" + po.PrePageText + "</a></li>"
	}
	return rs
}

/**
 * 尾部
 */
func getFooter(po *PageOptions, totalpages int) string {
	var rs = ""
	if po.EnablePreNexLink {
		var a = po.CurrentPage + 1
		if po.CurrentPage == totalpages {
			a = totalpages
		}
		rs += "<li " + judgeDisable(po, totalpages, 1) + "  ><a href='" + po.Href + "&" + po.ParamName + "=" + con.Itoa(a) + "'>" + po.NextPageText + "</a></li>"
	}
	if po.EnableFirstLastLink { //当首页,尾页都设定的时候,就显示
		rs += "<li " + judgeDisable(po, totalpages, 1) + " ><a href='" + po.Href + "&" + po.ParamName + "=" + con.Itoa(totalpages) + "'>" + po.LastPageText + "</a></li>"
	}
	rs += "</ul>"
	return rs
}

/**
 * 设置默认值
 */
func setDefault(po *PageOptions, totalpages int) *PageOptions {
	if len(po.FirstPageText) <= 0 {
		po.FirstPageText = "首页"
	}
	if len(po.LastPageText) <= 0 {
		po.LastPageText = "尾页"
	}
	if len(po.PrePageText) <= 0 {
		po.PrePageText = "上一页"
	}
	if len(po.NextPageText) <= 0 {
		po.NextPageText = "下一页"
	}
	if po.CurrentPage >= totalpages {
		po.CurrentPage = totalpages
	}
	if po.CurrentPage <= 1 {
		po.CurrentPage = 1
	}
	if po.LinkItemCount == 0 {
		po.LinkItemCount = 10
	}
	if po.PageSize == 0 {
		po.PageSize = 20
	}
	if len(po.ParamName) <= 0 {
		po.ParamName = "pno"
	}
	return po
}

/**
 *判断首页尾页  上一页下一页是否能用
 */
func judgeDisable(po *PageOptions, totalpages int, hOrf int) string {
	var rs = ""
	//判断头部
	if hOrf == 0 {
		if po.CurrentPage == 1 {
			rs = " "
		}
	} else {
		if po.CurrentPage == totalpages {
			rs = " "
		}
	}
	return rs
}
