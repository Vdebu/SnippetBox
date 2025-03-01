package mocks

import (
	"SnippetBox.mikudayo.net/internal/models"
	"time"
)

// 创建固定的snippet信息用于测试
var mockSnippet = &models.Snippet{
	ID:      1,
	Title:   "miku",
	Content: "mikudayo",
	Created: time.Now(),
	Expires: time.Now(),
}

// MockSnippetModel 不链接真实的数据库
type SnippetModel struct {
}

// 测试错误数据是否都正确返回

func (m *SnippetModel) Insert(title string, content string, expires int) (int, error) {
	return 2, nil
}

func (m *SnippetModel) Get(id int) (*models.Snippet, error) {
	switch id {
	case 1:
		return mockSnippet, nil
	default:
		return nil, models.ErrNoRecord
	}
}

func (m *SnippetModel) Latest() ([]*models.Snippet, error) {
	return []*models.Snippet{mockSnippet}, nil
}
