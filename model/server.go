package model

import (
	"github.com/zihang5127/easy-operation/common"
	"strconv"
	"time"

	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
)

// 服务器对象
type Server struct {
	Id         int       `orm:"pk;auto;unique;column(id)" json:"server_id"`
	Name       string    `orm:"size(255);column(name)" json:"name"`
	Type       string    `orm:"size(255);column(type);default(SSH)" json:"type"`
	IpAddress  string    `orm:"size(255);column(ip_address)" json:"ip_address"`
	Port       int       `orm:"type(int);column(port)" json:"port"`
	Username   string    `orm:"size(255);column(user_name)" json:"user_name"`
	PrivateKey string    `orm:"type(text);column(private_key)" json:"private_key"`
	Tag        string    `orm:"size(1000);column(tag)" json:"tag"`
	Status     int       `orm:"type(int);column(status);default(0)" json:"status"`
	CreateTime time.Time `orm:"type(datetime);column(create_time);auto_now_add" json:"create_time"`
	CreateAt   int       `orm:"type(int);column(create_at)" json:"-"`
}

// 获取对应数据库表名
func (m *Server) TableName() string {
	return "eo_server"
}

// 获取数据使用的引擎
func (m *Server) TableEngine() string {
	return "INNODB"
}

// 新建服务器对象
func NewServer() *Server {
	return &Server{}
}

// 根据ID查找对象
func (m *Server) Find() (error) {

	if m.Id <= 0 {
		return ErrInvalidParameter
	}
	o := orm.NewOrm()

	if err := o.Read(m); err != nil {
		logs.Error("%s: ", "v0.1")
		return err
	}
	return nil;
}

// 创建或更新
func (m *Server) Save() error {
	o := orm.NewOrm()
	var err error;
	if m.Id > 0 {
		_, err = o.Update(m)
	} else {
		_, err = o.Insert(m)
	}

	return err
}

// 删除
func (m *Server) Delete() error {
	o := orm.NewOrm()
	_, err := o.Delete(m)

	return err
}

// 搜索指定用户的服务器
func (m *Server) Search(keyword string, userId int, excludeServerId ...int) ([]Server, error) {
	o := orm.NewOrm()

	keyword = "%" + keyword + "%"

	sql := common.SearchServerByKeyword

	if len(excludeServerId) > 0 {
		sql += " AND id not in ("
		for _, num := range excludeServerId {
			sql += strconv.Itoa(num) + ","
		}
		sql += "0)"
	}

	var servers []Server

	_, err := o.Raw(sql, userId, keyword, keyword).QueryRows(&servers)

	if err != nil {
		logs.Error("", err.Error())
		return servers, err
	}

	return servers, nil
}

// 根据server_id和用户id查询服务器信息列表
func (m *Server) QueryServerByServerId(serverIds []int, userId ...int) ([]*Server, error) {
	o := orm.NewOrm()

	query := o.QueryTable(m.TableName()).Filter("id__in", serverIds)

	if len(userId) > 0 {
		query = query.Filter("create_at", userId[0])
	}

	var servers []*Server

	_, err := query.All(&servers)

	return servers, err
}
