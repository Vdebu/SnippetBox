package models

import (
	"database/sql"
	"time"
)

// 存储用户信息的结构体(与数据库中表的结构一致)
type User struct {
	ID             int
	Name           string
	Email          string
	HashedPassword []byte
	Created        time.Time
}

// 注入数据库依赖
type UserModel struct {
	DB *sql.DB
}

// 在数据库中新建用户
func (m *UserModel) Insert(name, email, password string) error {
	return nil
}

// 检查是否存在该用户 如果存在就返回id
func (m *UserModel) Authenticate(email, password string) (int, error) {
	return 0, nil
}

// 通过提供的id检查用户是否存在
func (m *UserModel) Exist(id int) (bool, error) {
	return false, nil
}
