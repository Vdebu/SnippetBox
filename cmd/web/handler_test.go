package main

import (
	"SnippetBox.mikudayo.net/internal/assert"
	"net/http"
	"testing"
)

func TestPing(t *testing.T) {
	// 定义用于后续测试的app实例 只初始化必要的底层依赖 否则会发生panic
	app := newTestApplication()
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
