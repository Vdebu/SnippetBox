package main

import (
	"SnippetBox.mikudayo.net/internal/assert"
	"net/http"
	"net/http/cookiejar"
	"net/url"
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

func TestUserSignup(t *testing.T) {
	// 创建后续用于测试的app
	app := newTestApplication(t)
	// 创建测试服
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	// 向指定处理器发送GET请求 并提取CSRF token尝试输出
	_, _, body := ts.get(t, "/user/signup")
	// 尝试提取token 若提取失败返回的会是空字符串
	validCSRFToken := extractCSRFToken(t, body)

	// 使用测试log输出提取到的csrfToken go test -v -run="TestUserSignup"
	//t.Logf("CSRF token is: %q", csrfToken)

	// 使用表驱动测试
	const (
		validName     = "miku"
		validPassword = "mikudayo3939"
		validEmail    = "miku@vocaloid.com"
		// 在这里会对body中的字符进行严格匹配 精确到字符大小写 引号以及空格！
		// 确认页面中的表单元素是否正确生成 从而保证用户填写数据后能够正确提交和被服务器解析
		formTag = "<form action='/user/signup' method = 'POST' novalidate>"
	)
	tests := []struct {
		name         string
		userName     string
		userEmail    string
		userPassword string
		csrfToken    string
		wantCode     int
		wantFormTag  string
	}{
		{
			name:         "Valid submission",
			userName:     validName,
			userEmail:    validEmail,
			userPassword: validPassword,
			csrfToken:    validCSRFToken,
			wantCode:     http.StatusSeeOther,
		},
		{
			name:         "Invalid CSRF token",
			userName:     validName,
			userEmail:    validEmail,
			userPassword: validPassword,
			csrfToken:    "wrongToken",
			wantCode:     http.StatusBadRequest,
		},
		{
			name:         "Empty name",
			userName:     "",
			userEmail:    validEmail,
			userPassword: validPassword,
			csrfToken:    validCSRFToken,
			wantCode:     http.StatusUnprocessableEntity,
			wantFormTag:  formTag,
		},
		{
			name:         "Empty email",
			userName:     validName,
			userEmail:    "",
			userPassword: validPassword,
			csrfToken:    validCSRFToken,
			wantCode:     http.StatusUnprocessableEntity,
			wantFormTag:  formTag,
		},
		{
			name:         "Empty password",
			userName:     validName,
			userEmail:    validEmail,
			userPassword: "",
			csrfToken:    validCSRFToken,
			wantCode:     http.StatusUnprocessableEntity,
			wantFormTag:  formTag,
		},
		{
			name:         "Invalid email",
			userName:     validName,
			userEmail:    "pa$$",
			userPassword: validPassword,
			csrfToken:    validCSRFToken,
			wantCode:     http.StatusUnprocessableEntity,
			wantFormTag:  formTag,
		},
		{
			name:         "Short password",
			userName:     validName,
			userEmail:    validEmail,
			userPassword: "pa$$",
			csrfToken:    validCSRFToken,
			wantCode:     http.StatusUnprocessableEntity,
			wantFormTag:  formTag,
		},
		{
			name:         "Duplicate email",
			userName:     validName,
			userEmail:    "teto@vocaloid.com",
			userPassword: validPassword,
			csrfToken:    validCSRFToken,
			wantCode:     http.StatusUnprocessableEntity,
			wantFormTag:  formTag,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			form := url.Values{}
			// 初始化url表单数据
			// 这些都是signup网页中的html标签name 以键值的形式加入模拟填写用于后续解析
			form.Add("name", tt.userName)
			form.Add("email", tt.userEmail)
			form.Add("password", tt.userPassword)
			form.Add("csrf_token", tt.csrfToken)

			// 尝试发送Post请求并获取状态码与响应体
			code, _, body := ts.postForm(t, "/user/signup", form)
			// 判断是否如预期
			assert.Equal(t, code, tt.wantCode)
			// 如果有formTag就进行测试
			if tt.wantFormTag != "" {
				assert.StringContains(t, body, tt.wantFormTag)
			}
		})
	}
}

func TestSnippetCreate(t *testing.T) {
	// 创建用于测试的服务器与Application
	app := newTestApplication(t)
	ts := newTestServer(t, app.routes())
	// 使用完毕关闭服务器
	defer ts.Close()
	_, _, body := ts.get(t, "/user/login")
	csrfToken := extractCSRFToken(t, body)
	const (
		validEmail    = "miku@vocaloid.com"
		validPassword = "mikudayo3939"
		wantLocation  = "/user/login"
		wantForm      = "<form action='/snippet/create' method='POST'>"
	)
	// 测试执行的顺序也很重要 因为登入成功后对cookie的改变会影响后续的测试结果
	t.Run("Unauthentivated", func(t *testing.T) {
		// 在未登入的状态下直接对/snippet/create发送访问请求 结果会被重定向
		code, header, _ := ts.get(t, "/snippet/create")
		assert.Equal(t, code, http.StatusSeeOther)
		assert.Equal(t, header.Get("Location"), wantLocation)
	})
	t.Run("Authenticated", func(t *testing.T) {
		form := url.Values{}
		form.Add("email", validEmail)
		form.Add("password", validPassword)
		form.Add("csrf_token", csrfToken)
		ts.postForm(t, "/user/login", form)
		code, _, body := ts.get(t, "/snippet/create")
		// 判断是否成功进入网页
		assert.Equal(t, code, http.StatusOK)
		assert.StringContains(t, body, wantForm)
	})

}

// 神秘原因 解决了后续无效信息访问请求的问题 但是有效信息却无法正确登入了 疑似cookie过度清除
func TestSC(t *testing.T) {
	app := newTestApplication(t)
	ts := newTestServer(t, app.routes())
	// 使用完毕关闭服务器
	defer ts.Close()
	const (
		validEmail    = "miku@vocaloid.com"
		validPassword = "mikudayo3939"
		wantLocation  = "/user/login"
		wantForm      = "<form action='/snippet/create' method='POST'>"
	)
	//在这里只需要测试非登录状态与登录状态是否能成功访问创建消息的页面即可
	//不需要进行额外的测试如判断是否是合法的登入信息 一个函数只做一件事
	tests := []struct {
		name         string
		email        string
		password     string
		wantCode     int
		wantLocation string
		wantForm     string
	}{
		{
			name:         "Valid user info",
			email:        validEmail,
			password:     validPassword,
			wantCode:     http.StatusOK,
			wantForm:     wantForm,
			wantLocation: "/snippet/create",
		},
		{
			name:         "Empty user email info",
			email:        "",
			password:     validPassword,
			wantCode:     http.StatusSeeOther,
			wantLocation: wantLocation,
		},
		{
			name:         "Empty user password info",
			email:        validEmail,
			password:     "",
			wantCode:     http.StatusSeeOther,
			wantLocation: wantLocation,
		},
		{
			name:         "Empty user password and email info",
			email:        "",
			password:     "",
			wantCode:     http.StatusSeeOther,
			wantLocation: wantLocation,
		},
	}
	// 确保每个测试子案例使用独立的会话!!!
	// 在每个子测试开始前显式清除或重置 session，以确保无效登录测试时不会受到之前成功登录的影响
	for _, tt := range tests {
		// 载入测试名称信息与测试用函数
		t.Run(tt.name, func(t *testing.T) {
			// 每一次测试显示清除原先留下的cookie
			jar, err := cookiejar.New(nil)
			if err != nil {
				t.Fatal(err)
			}
			ts.Client().Jar = jar
			_, _, body := ts.get(t, "/user/login")
			csrfToken := extractCSRFToken(t, body)
			// 先登入再向创建消息发送请求
			form := url.Values{}
			// 模拟用户填写登入信息
			form.Add("email", tt.email)
			form.Add("password", tt.password)
			form.Add("csrf_token", csrfToken)
			// 将填写好的信息传入postForm用于登入
			ts.postForm(t, "/user/login", form)
			// 登入后尝试向/snippet/create发送请求te
			code, header, _ := ts.get(t, "/snippet/create")
			// 判断值是否如预期
			assert.Equal(t, code, tt.wantCode)
			// 检查是否跳转到登录页面
			assert.Equal(t, header.Get("Location"), tt.wantLocation)
			//如果无效登录的情况下并没有进行 logout（因为登录处理器直接渲染登录页面），session 里的旧数据可能依然存在
			// 采用直接清除诜cookie的方式进行处理
			//ts.postForm(t, "/user/logout", nil)
		})
	}
}
