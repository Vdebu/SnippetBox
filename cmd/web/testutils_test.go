package main

import (
	"SnippetBox.mikudayo.net/internal/models/mocks"
	"bytes"
	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
	"time"

	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"testing"
)

// 生成用于测试的Application 只初始化必要的底层依赖
func newTestApplication(t *testing.T) *Application {
	// 创建所有必要的依赖
	// 创建网页模板的缓存
	templateCache, err := newTemplateCache()
	if err != nil {
		t.Fatal(err)
	}
	// 创建解码器
	formDecoder := form.NewDecoder()
	// 创建session
	sessionManager := scs.New()
	sessionManager.Lifetime = 12 * time.Hour
	sessionManager.Cookie.Secure = true
	return &Application{
		// app中各处的方法都用到了自定义log 不初始化会发生panic
		errlog:         log.New(io.Discard, "", 0),
		infolog:        log.New(io.Discard, "", 0),
		templateCache:  templateCache,
		formDecoder:    formDecoder,
		sessionManager: sessionManager,
		snippets:       &mocks.SnippetModel{},
		users:          &mocks.UserModel{},
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
	// 初始化cookiejar用于服务器cookie的自动存储
	jar, err := cookiejar.New(nil)
	if err != nil {
		t.Fatal()
	}
	// 所有的cookie会被自动存储并在子请求中使用(通过当前客户端发送的子请求)
	ts.Client().Jar = jar
	// 关闭服务器的自动重定向
	ts.Client().CheckRedirect = func(req *http.Request, via []*http.Request) error {
		// 返回这个错误的同时最近一次的请求也会被返回,接下的请求会被终止
		// 在遇到比如 HTTP 3xx 状态码时会调用这个自定义的CheckRedirect函数
		return http.ErrUseLastResponse
	}
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
