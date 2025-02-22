package models

import (
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
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
	// 从用户输入的密码生成哈希 使用2^12(4096)次迭代
	// 这里哈希值的返回形式是字节 后续向数据库中进行插入要进行字符串的转换
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}
	// 尝试向数据库中插入新用户
	stmt := `INSERT INTO users(name,email,hashed_password,created)
	VALUES(?,?,?,UTC_TIMESTAMP())`
	_, err = m.DB.Exec(stmt, name, email, string(hashedPassword))
	if err != nil {
		// 对sql的报错进行特判
		// 像先前特判从网页解码数据一样使用errors.AS()进行判断
		var mySQLError *mysql.MySQLError
		if errors.As(err, &mySQLError) {
			// 如果是sql报的错
			if mySQLError.Number == 1062 && strings.Contains(mySQLError.Message, "users_uc_email") {
				// 并且错误代码与索引匹配
				// 返回自定义错误
				return ErrDuplicateEmail
			}
		}
		return err
	}
	return nil
}

// 检查是否存在该用户 如果存在就返回id
func (m *UserModel) Authenticate(email, password string) (int, error) {
	// 定义变量用于从数据库中提取数据
	var id int
	var hashedPassword []byte

	stmt := `SELECT id,hashed_password FROM users WHERE email = ?`
	err := m.DB.QueryRow(stmt, email).Scan(&id, &hashedPassword)
	if err != nil {
		// 判断是否为sql查询为空的错误
		if errors.Is(err, sql.ErrNoRows) {
			// 如果查询结果为空值直接返回无效数据错误
			return 0, ErrInvalidCredentials
		} else {
			// 其他错误统一正常返回处理
			return 0, err
		}
	}
	// 确实存在这个邮箱 检查用户填写的密码哈希值与数据库中存储的是否一致
	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err != nil {
		// 判断是否是哈希值不匹配的错误
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			// 如果哈希值不匹配直接返回无效数据错误
			return 0, ErrInvalidCredentials
		} else {
			return 0, err
		}
	}
	// 登陆成功
	return id, nil
}

// 通过提供的id检查用户是否存在
func (m *UserModel) Exists(id int) (bool, error) {
	return false, nil
}
