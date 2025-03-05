package mocks

import (
	"SnippetBox.mikudayo.net/internal/models"
)

type UserModel struct {
}

// 测试错误数据是否都正确返回

func (m *UserModel) Insert(name, email, password string) error {
	// 在后续测试会进行调用 用于判断是否使用了重复的邮箱
	switch email {
	case "teto@vocaloid.com":
		return models.ErrDuplicateEmail
	default:
		return nil
	}
}
func (m *UserModel) Authenticate(email, password string) (int, error) {
	if email == "miku@vocaloid.com" && password == "mikudayo3939" {
		return 39, nil
	}
	return 0, models.ErrInvalidCredentials
}
func (m *UserModel) Exists(id int) (bool, error) {
	switch id {
	case 39:
		return true, nil
	default:
		return false, nil
	}
}

// 返回用户的账号名
func (m *UserModel) GetName(id int) (string, error) {
	return "", nil
}

// 返回用户账号的创建时间
func (m *UserModel) GetJoinedTime(id int) (string, error) {
	return "", nil
}

// 返回用户的邮箱
func (m *UserModel) GetEmail(id int) (string, error) {

	return "", nil

}
