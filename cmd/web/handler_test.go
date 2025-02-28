package main

import (
	"SnippetBox.mikudayo.net/internal/assert"
	"bytes"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPing(t *testing.T) {
	// 定义用于后续测试的app实例 只初始化必要的底层依赖 否则会发生panic
	app := &Application{
		infolog: log.New(io.Discard, "", 0),
		errlog:  log.New(io.Discard, "", 0),
	}
	// 创建一个HTTPS协议的测试服务器 将定义在路由模块中的路由全部传入用服务器请求的处理
	// 这个测试服务器会随机监听一个本地端口
	// 如果是HTTP服务器应该使用httptest.NewServer()
	ts := httptest.NewTLSServer(app.routes())
	// 在测试结束之后将服务器关闭
	defer ts.Close()
	// 尝试向服务器的指定路由发送指定请求
	// ts.URL存储了整个服务器的前缀 链接上要访问的路径即可
	// https://localhost:xxxx/
	rs, err := ts.Client().Get(ts.URL + "/ping")
	if err != nil {
		t.Fatal(err)
	}
	// 检查状态码是如预期
	assert.Equal(t, rs.StatusCode, http.StatusOK)

	// 检查响应体中的内容是否如预期
	defer rs.Body.Close()
	body, err := io.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}
	// 去除空格后进行比较
	bytes.TrimSpace(body)
	assert.Equal(t, string(body), "OK")
}
