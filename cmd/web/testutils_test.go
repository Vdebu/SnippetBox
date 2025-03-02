package main

import (
	"SnippetBox.mikudayo.net/internal/models/mocks"
	"bytes"
	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
	"html"
	"net/url"
	"regexp"
	"time"

	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"testing"
)

// 使用正则表达式对csrfToken进行匹配 (.+)表示重复匹配引号中的任意字符
var csrfTokenRX = regexp.MustCompile(`<input type="hidden" name="csrf_token" value="(.+)">`)

func extractCSRFToken(t *testing.T, body string) string {
	matches := csrfTokenRX.FindStringSubmatch(body)
	if len(matches) < 2 {
		t.Fatal("no csrf token found in body...")
	}
	// 将切片进行转换并返回
	return html.UnescapeString(string(matches[1]))
}

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

// 在测试服上绑定一个post方法接收指定的url(目标路由)与要发送的目标值 返回状态码,响应头,响应体
func (ts *testServer) postForm(t *testing.T, urlPath string, form url.Values) (int, http.Header, string) {
	// 向指定的路由发送post请求
	// 将传入的数据编码为application/x-www-form-urlencoded格式,并以POST请求的方式发送到指定的URL
	// 模拟了用户在网页上填写表单后点击提交时浏览器所执行的操作
	rs, err := ts.Client().PostForm(ts.URL+urlPath, form)
	if err != nil {
		t.Fatal(err)
	}
	// 之后定义在处理器中的FormDecoder会自动将数据进行解析并用于后续使用
	// 使用完毕关闭响应体
	defer rs.Body.Close()
	body, err := io.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}
	bytes.TrimSpace(body)
	// 返回响应体的信息
	return rs.StatusCode, rs.Header, string(body)
}
