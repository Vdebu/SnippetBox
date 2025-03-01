package main

import (
	"SnippetBox.mikudayo.net/internal/assert"
	"net/http"
	"testing"
)

func TestPing(t *testing.T) {
	// 定义用于后续测试的app实例 只初始化必要的底层依赖 否则会发生panic
	app := newTestApplication(t)
	// 创建一个HTTPS协议的测试服务器 将定义在路由模块中的路由全部传入用服务器请求的处理
	// 这个测试服务器会随机监听一个本地端口
	// 如果是HTTP服务器应该使用httptest.NewServer()
	ts := newTestServer(t, app.routes())
	// 在测试结束之后将服务器关闭
	defer ts.Close()
	// 尝试向服务器的指定路由发送指定请求
	// ts.URL存储了整个服务器的前缀 链接上要访问的路径即可
	// https://localhost:xxxx/
	code, _, body := ts.get(t, "/ping")
	// 检查返回值是否如预期
	assert.Equal(t, code, http.StatusOK)
	assert.Equal(t, body, "OK")
}

func TestSnippetView(t *testing.T) {
	// 创建mock app用于后续测试
	app := newTestApplication(t)
	// 创建测试服务器
	ts := newTestServer(t, app.routes())
	// 使用完毕后关闭服务器
	defer ts.Close()
	// 使用表驱动对处理器进行测试
	tests := []struct {
		name     string
		urlPath  string
		wantCode int
		wantBody string
	}{
		{
			name:     "Valid ID",
			urlPath:  "/snippet/view/39",
			wantCode: http.StatusOK,
			wantBody: "mikudayo",
		},
		{
			name:     "Non-existent ID",
			urlPath:  "/snippet/view/93",
			wantCode: http.StatusNotFound,
		},
		{
			name:     "Negative ID",
			urlPath:  "/snippet/view/-39",
			wantCode: http.StatusNotFound,
		},
		{
			name:     "Decimal ID",
			urlPath:  "/snippet/view/3.39",
			wantCode: http.StatusNotFound,
		},
		{
			name:     "String ID",
			urlPath:  "/snippet/view/miku",
			wantCode: http.StatusNotFound,
		},
		{
			name:     "Empty ID",
			urlPath:  "/snippet/view/",
			wantCode: http.StatusNotFound,
		},
	}
	for _, tt := range tests {
		// 启动子测试
		t.Run(tt.name, func(t *testing.T) {
			// 使用自定义类型的绑定方法
			code, _, body := ts.get(t, tt.urlPath)
			// 检查HTTP状态码
			assert.Equal(t, code, tt.wantCode)
			// 有些测试用例是不需要检查响应体的即wantBody == "" 进行特判处理
			if tt.wantBody != "" {
				// 只需要检查响应体中是否有特定的字符串即可
				assert.StringContains(t, body, tt.wantBody)
			}
		})
	}
}
