package model

import (
	"github.com/zihang5127/easy-operation/module/encry"
	"time"

	"github.com/astaxie/beego/orm"
)

type User struct {
	Id            int       `orm:"pk;auto;unique;column(id)"`
	Username      string    `orm:"size(100);unique;column(user_name)"`
	Password      string    `orm:"size(1000);column(password)"`
	Email         string    `orm:"size(255);column(email);null;default(null)"`
	Phone         string    `orm:"size(255);column(phone);null;default(null)"`
	Avatar        string    `orm:"size(1000);column(avatar)"`
	Role          int       `orm:"column(role);type(int);default(1)"`
	Status        int       `orm:"column(status);type(int);default(0)"`
	CreateTime    time.Time `orm:"type(datetime);column(create_time);auto_now_add"`
	UpdateTime    time.Time `orm:"type(datetime);column(update_time);null"`
	LastLoginTime time.Time `orm:"type(datetime);column(last_login_time);null"`
}

func (u *User) TableName() string {
	return "eo_user"
}

/**
 * 获取数据库引擎
 */
func (u *User) TableEngine() string {
	return "INNODB"
}

// 新建用户对象
func NewUser() *User {
	return new(User)
}

// 根据id查找用户
func (u *User) Find() error {
	o := orm.NewOrm()
	err := o.Read(u)

	if err == orm.ErrNoRows {
		return ErrUserNoExist
	}

	return nil
}

// 登录
func (u *User) Login(username string, password string) (*User, error) {
	o := orm.NewOrm()

	user := &User{}

	err := o.QueryTable(u.TableName()).Filter("username", username).Filter("status", 0).One(user)

	if err != nil {
		return user, ErrUserNoExist
	}

	ok, err := encry.PasswordVerify(user.Password, password)

	if ok && err == nil {
		return user, nil
	}

	return user, ErrorUserPasswordError
}

// 添加用户
func (u *User) Add() error {
	o := orm.NewOrm()

	hash, err := encry.PasswordHash(u.Password)

	if err != nil {
		return err
	}

	u.Password = hash

	_, err = o.Insert(u)

	if err != nil {
		return err
	}
	return nil
}

// 更新
func (u *User) Update(cols ... string) (error) {
	o := orm.NewOrm()

	if _, err := o.Update(u, cols...); err != nil {
		return err
	}
	return nil
}

// 删除
func (u *User) Delete() error {
	o := orm.NewOrm()

	if _, err := o.Delete(u); err != nil {
		return err
	}
	return nil
}
