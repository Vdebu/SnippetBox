package main

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

// 生成用于测试的Application 只初始化必要的底层依赖
func newTestApplication() *Application {
	return &Application{
		// app中各处的方法都用到了自定义log 不初始化会发生panic
		errlog:  log.New(io.Discard, "", 0),
		infolog: log.New(io.Discard, "", 0),
	}
}

// 使用结构体嵌入测试会用的服务器
type testServer struct {
	*httptest.Server
}

// 创建并返回一个测试服务器实例
func newTestServer(t *testing.T, h http.Handler) *testServer {
	// 以传入的路由为基础建立测试服务器
	ts := httptest.NewTLSServer(h)
	return &testServer{ts}
}

// 在测试服上绑定一个GET方法接收指定的url(目标路由) 返回状态码,响应头,响应体
func (ts *testServer) get(t *testing.T, urlPath string) (int, http.Header, string) {
	// 向测试服发送GET请求
	rs, err := ts.Client().Get(ts.URL + urlPath)
	if err != nil {
		t.Fatal(err)
	}
	// 读取完毕响应体后进行关闭
	defer rs.Body.Close()
	body, err := io.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}
	bytes.TrimSpace(body)
	// 返回处理完毕的结果
	return rs.StatusCode, rs.Header, string(body)
}
