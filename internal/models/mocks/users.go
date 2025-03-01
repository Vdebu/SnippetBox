package mocks

import "SnippetBox.mikudayo.net/internal/models"

type UserModel struct {
}

// 测试错误数据是否都正确返回

func (m *UserModel) Insert(name, email, password string) error {
	switch email {
	case "mikudayo@vocaloid.com":
		return models.ErrDuplicateEmail
	default:
		return nil
	}
}
func (m *UserModel) Authenticate(email, password string) (int, error) {
	if email == "mikudayo@vocaloid.com" && password == "mikudayo39" {
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
