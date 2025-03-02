package models

import (
	"SnippetBox.mikudayo.net/internal/assert"
	"testing"
)

func TestUserModelExists(t *testing.T) {
	if testing.Short() {
		// 若在进行测试时提供了-short flag就会跳过TestUserModelExists go test -v -short ./...
		t.Skip("models:skipping integration test...")
	}
	tests := []struct {
		name   string
		userID int
		want   bool
	}{
		{
			name:   "Valid ID",
			userID: 39,
			want:   true,
		},
		{
			name:   "Zero ID",
			userID: 0,
			want:   false,
		},
		{
			name:   "Non-existent ID",
			userID: 93,
			want:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 创建测试用数据库 每次数据库中的资源由于先前的设计逻辑都会自动refresh
			db := newTestDB(t)
			// 实例化用户模型 这里不用mock 直接链接数据库进行操作
			m := UserModel{DB: db}
			// 调用Exists方法检查预期值
			exists, err := m.Exists(tt.userID)
			// 判断布尔值是否如预期
			assert.Equal(t, exists, tt.want)
			// 判断是否未发生错误
			assert.NilError(t, err)
		})
	}
}
