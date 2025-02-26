package main

import (
	"testing"
	"time"

	"SnippetBox.mikudayo.net/internal/assert"
)

// 为模板中的函数编写测试

func TestHumanDate(t *testing.T) {
	// 使用匿名结构体存储测试会用上的数据(测试用例与预期结果)
	tests := []struct {
		name string
		tm   time.Time
		want string
	}{
		{
			name: "UTC",
			tm:   time.Date(2022, 3, 7, 10, 15, 0, 0, time.UTC),
			want: "2022-03-07 10:15:00",
		},
		{
			name: "Empty",
			tm:   time.Time{},
			want: "",
		},
		{
			name: "CET",
			// 延迟一个小时
			tm:   time.Date(2022, 3, 7, 10, 15, 0, 0, time.FixedZone("CET", 1*60*60)),
			want: "2022-03-07 09:15:00",
		},
	}

	// 通过遍历测试用例来测试程序
	for _, tt := range tests {
		// 相当于是启动了多个goroutine来测试程序
		t.Run(tt.name, func(t *testing.T) {
			hd := hunmanDate(tt.tm)
			// 使用自定义泛型函数进行比较
			assert.Equal(t, hd, tt.want)
		})
	}
	// 创建time.Time对象传入函数
	// tm := time.Date(2022, 3, 17, 10, 15, 0, 0, time.UTC)
	// hd := hunmanDate(tm)
	// 检查输出是否如预期
	// if hd != "2022-03-17 10:15:00" {
	// 	t.Errorf("got %q;want %q", hd, "2022-03-17 10:15:00")
	// }
}
