package main

import (
	"net/http"

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
	fs := http.FileServer(http.Dir("D:/Program/Mycode/Now/Mygo/Project/main/SnippetBox/ui/static"))
	// 去除前缀后从文件服务器中查找文件并返回
	// 不想让用户直接访问根目录可以检测访问路径并直接返回一个静态页面
	router.Handler(http.MethodGet, "/static/*filepath", http.StripPrefix("/static", fs))

	// 创建包含seesion的新中间件链对需要共享信息的路由进行手动预包装
	dynamic := alice.New(app.sessionManager.LoadAndSave)

	// .ThenFunc()返回的还是一个handler而不是像HandlerFunc直接成为可执行的路由 所以在这里要改变原先router.HandlerFunc()为router.Handler()来注册路由
	router.Handler(http.MethodGet, "/", dynamic.ThenFunc(app.home))
	router.Handler(http.MethodGet, "/snippet/view/:id", dynamic.ThenFunc(app.snippetView))
	// render 一个用于填写信息的网页
	router.Handler(http.MethodGet, "/snippet/create", dynamic.ThenFunc(app.snippetCreate))
	router.Handler(http.MethodPost, "/snippet/create", dynamic.ThenFunc(app.snippetCreatePost))

	// 使用包创建一个中间件链变量方便管理 执行顺序 ->
	standard := alice.New(app.recoverPanic, app.logRequest, secureHeaders)
	// 使用中间件将当前mux下的所有路由都包装起来
	// 相当于是"重写"的在结构体中的方法
	// 最外层的中间件会第一个进行应用 类似于栈 first in first out
	// return app.recoverPanic(app.logRequest(secureHeaders(mux)))

	// 直接调用then方法初始化路由
	return standard.Then(router)
}
