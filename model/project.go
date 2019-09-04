package model

import (
	"errors"
	"github.com/astaxie/beego/orm"
	"github.com/zihang5127/easy-operation/module/encry"
	"time"
)

// User 会员信息.
type Project struct {
	Id             int       `orm:"pk;auto;unique;column(id)" json:"id"`
	RepositoryName string    `orm:"size(255);column(repo_name)" json:"repository_name"`
	BranchName     string    `orm:"size(255);column(branch_name)" json:"branch_name"`
	Tag            string    `orm:"size(1000);column(tag)" json:"tag"`
	Shell          string    `orm:"size(1000);column(shell)" json:"shell"`
	Status         int       `orm:"type(int);column(status);default(0)" json:"status"`
	Key            string    `orm:"size(255);column(key);unique" json:"key"`
	Secure         string    `orm:"size(255);column(secure)" json:"secure"`
	LastExecTime   time.Time `orm:"type(datetime);column(last_exec_time);null" json:"last_exec_time"`
	RepositoryType string    `orm:"size(20);column(repo_type)" json:"repository_type"`
	CreateAt       int       `orm:"type(int);column(create_at)"`
	CreateTime     time.Time `orm:"type(datetime);column(create_time);auto_now_add" json:"create_time"`
	UpdateTime     time.Time `orm:"type(datetime);column(update_time);auto_now_add" json:"update_time"`
}

func (m *Project) TableName() string {
	return "eo_project"
}

func (m *Project) TableEngine() string {
	return "INNODB"
}

func NewProject() *Project {
	return &Project{}
}

func (m *Project) Find() error {

	if m.Id <= 0 {
		return ErrInvalidParameter
	}

	o := orm.NewOrm()

	if err := o.Read(m); err != nil {
		return err
	}
	return nil;
}

// 批量删除
func (m *Project) DeleteMulti(id ...int) error {
	if len(id) > 0 {
		o := orm.NewOrm()
		ids := make([]int, len(id))
		params := ""

		for i := 0; i < len(id); i++ {
			ids[i] = id[i]
			params = params + ",?"
		}
		_, err := o.Raw("DELETE FROM eo_project WHERE id IN ("+params[1:]+")", ids).Exec()
		if err != nil {
			return err
		}
		return nil
	}
	return errors.New("Invalid parameter")
}

//删除一条
func (m *Project) Delete() error {
	o := orm.NewOrm()
	_, err := o.Delete(m)

	return err
}

// 根据Key查找
func (m *Project) FindByKey(key string) error {
	o := orm.NewOrm()

	if err := o.QueryTable(m.TableName()).Filter("key", key).One(m); err != nil {
		return err
	}
	return nil
}

// 添加或更新
func (m *Project) Save() error {
	o := orm.NewOrm()
	var err error;
	if m.Id > 0 {
		_, err = o.Update(m)
	} else {
		key := (time.Now().String() + m.RepositoryName + m.BranchName)

		m.Key = encry.Md5(key)

		m.Secure = encry.Md5(key + key + time.Now().String())

		_, err = o.Insert(m)
	}

	return err
}
