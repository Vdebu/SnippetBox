package main

import (
	"net/http"

	"SnippetBox.mikudayo.net/ui"
	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

func (app *Application) routes() http.Handler {
	// 创建一个自定义路由
	// mux := http.NewServeMux()
	// 使用三方路由库建立一个可以制定处理器访问方法与url占位符的复用器
	router := httprouter.New()

	// 重写当前路由的内置notfound函数 使整个应用程序表现一致
	// 尝试访问不存在的路由器与合法但是不存在的页面
	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.notFound(w)
	})

	// 调用
	// 创建静态文件服务器
	// fs := http.FileServer(http.Dir("D:/Program/Mycode/Now/Mygo/Project/main/SnippetBox/ui/static"))

	// 使用嵌入的文件系统(embed.FS)作为FileServer
	fs := http.FileServer(http.FS(ui.Files))
	// 去除前缀后从文件服务器中查找文件并返回
	// 不想让用户直接访问根目录可以检测访问路径并直接返回一个静态页面
	// router.Handler(http.MethodGet, "/static/*filepath", http.StripPrefix("/static", fs))

	// 不需要再去除url前缀 直接传入即可
	router.Handler(http.MethodGet, "/static/*filepath", fs)

	// 创建用于测试的路由
	router.HandlerFunc(http.MethodGet, "/ping", ping)

	// 创建包含seesion的新中间件链对需要共享信息的路由进行手动预包装
	// 添加防止CSRF攻击的noSurf中间件
	dynamic := alice.New(app.sessionManager.LoadAndSave, noSurf, app.authenticate)

	// .ThenFunc()返回的还是一个handler而不是像HandlerFunc直接成为可执行的路由 所以在这里要改变原先router.HandlerFunc()为router.Handler()来注册路由
	router.Handler(http.MethodGet, "/", dynamic.ThenFunc(app.home))
	router.Handler(http.MethodGet, "/snippet/view/:id", dynamic.ThenFunc(app.snippetView))
	// 用户信息处理相关的处理器
	router.Handler(http.MethodGet, "/user/signup", dynamic.ThenFunc(app.userSignup))
	router.Handler(http.MethodPost, "/user/signup", dynamic.ThenFunc(app.userSignupPost))
	// 用户登入相关的处理器
	router.Handler(http.MethodGet, "/user/login", dynamic.ThenFunc(app.userLogin))
	router.Handler(http.MethodPost, "/user/login", dynamic.ThenFunc(app.userLoginPost))

	// 对路由进行分组处理 上半部分的网页访问不需要用户的登入权限 在下半部分进行检测
	// 下面还需要合适用户的身份信息就用新的中间件 不会再次从数据库进行查询 直接从ctx中进行核实
	protected := dynamic.Append(app.requireAuthentication)

	// 创建消息相关的处理器
	router.Handler(http.MethodGet, "/snippet/create", protected.ThenFunc(app.snippetCreate))
	router.Handler(http.MethodPost, "/snippet/create", protected.ThenFunc(app.snippetCreatePost))
	// 用户退出的相关处理器
	router.Handler(http.MethodPost, "/user/logout", protected.ThenFunc(app.userLogoutPost))
	// 使用包创建一个中间件链变量方便管理 执行顺序 ->
	standard := alice.New(app.recoverPanic, app.logRequest, secureHeaders)
	// 使用中间件将当前mux下的所有路由都包装起来
	// 相当于是"重写"的在结构体中的方法
	// 最外层的中间件会第一个进行应用 类似于栈 first in first out
	// return app.recoverPanic(app.logRequest(secureHeaders(mux)))

	// 直接调用then方法初始化路由
	return standard.Then(router)
}
