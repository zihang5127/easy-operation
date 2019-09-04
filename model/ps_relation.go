package model

import (
	"errors"
	"github.com/zihang5127/easy-operation/common"
	"time"

	"github.com/astaxie/beego/orm"
)

// PsRelation Project和Server之间关系
type PsRelation struct {
	Id         int       `orm:"pk;auto;unique;column(id)" json:"id"`
	ProjectId  int       `orm:"type(int);column(project_id)" json:"project_id"`
	ServerId   int       `orm:"type(int);column(server_id)" json:"server_id"`
	UserId     int       `orm:"type(int);column(user_id)" json:"user_id"`
	CreateTime time.Time `orm:"type(datetime);column(create_time);auto_now_add" json:"create_time"` //添加时间
}

// TableName 获取对应数据库表名
func (m *PsRelation) TableName() string {
	return "eo_ps_relation"
}

// TableEngine 获取数据使用的引擎
func (m *PsRelation) TableEngine() string {
	return "INNODB"
}

// Relation 获取新的关系对象
func NewPsRelation() *PsRelation {
	return &PsRelation{}
}

// Save 更新或添加映射关系
func (m *PsRelation) Save() error {

	o := orm.NewOrm()
	if o.QueryTable(m.TableName()).Filter("project_id", m.ProjectId).Filter("server_id", m.ServerId).Exist() {
		return ErrServerAlreadyExist
	}
	var err error

	if m.Id > 0 {
		if m.ProjectId <= 0 || m.ServerId <= 0 {
			return errors.New("Data format error")
		}
		_, err = o.Update(m)
	} else {
		_, err = o.Insert(m)
	}
	return err
}

// Delete 删除关系
func (m *PsRelation) Delete() error {
	o := orm.NewOrm()
	_, err := o.Delete(m)

	return err
}

// Find 查找关系
func (m *PsRelation) Find(id int) error {
	o := orm.NewOrm()

	m.Id = id

	if err := o.Read(m); err != nil {
		return err
	}
	return nil;
}

// RelationDetailed 包含 Proje 和 Server 信息的关系实体
type RelationDetailed struct {
	Id             int    `json:"id"`
	UserId         int    `json:"user_id"`
	ProjectId      int    `json:"project_id" orm:"column(project_id)"`
	Repository     string `json:"repository"`
	Branch         string `json:"branch"`
	ServerId       int    `json:"server_id"`
	ProjectTag     string `json:"project_tag"`
	Shell          string `json:"shell"`
	ProjectStatus  int    `json:"project_status"`
	RepositoryType string `json:"repository_type"`

	ServerName   string `json:"server_name"`
	ServerType   string `json:"server_type"`
	IpAddress    string `json:"ip_address"`
	Port         int    `json:"port"`
	Username     string `json:"username"`
	PrivateKey   string `json:"-"`
	ServerTag    string `json:"server_tag"`
	ServerStatus int    `json:"server_status"`
}

// FindRelationDetailedByWhere 指定条件查询完整的关系对象
func FindRelationDetailedByWhere(where string, params ...interface{}) ([]RelationDetailed, error) {
	o := orm.NewOrm()

	sql := common.FindRelationDetailedByWhereSql;
	if where != "" {
		sql += where
	}

	rawSetter := o.Raw(sql, params)

	var results []RelationDetailed

	_, err := rawSetter.QueryRows(&results)

	return results, err
}

// ServerRelation 服务与Project简单关系
type ServerRelation struct {
	ServerId   int
	Id         int
	ProjectId  int
	UserId     int
	Status     int
	Name       string
	IpAddress  string
	Port       int
	Type       string
	CreateTime time.Time
	CreateAt   int
}

// QueryByProjectId 查找指定用户的服务和Project简单关系
func (m *PsRelation) QueryByProjetId(projectId int, userId int) ([]*ServerRelation, error) {
	o := orm.NewOrm()

	var res []*ServerRelation

	sql := common.QueryByProjetIdSql

	_, err := o.Raw(sql, projectId, userId).QueryRows(&res)

	return res, err
}

// DeleteByWhere 删除指定用户的服务和Project的关系
func (m *PsRelation) DeleteByWhere(where string, args ...interface{}) error {
	o := orm.NewOrm()

	sql := "DELETE FROM eo_ps_relation WHERE 1=1 " + where

	_, err := o.Raw(sql, args).Exec()

	return err
}
