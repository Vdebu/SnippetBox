package main

import "net/http"

func (app *Application) routes() http.Handler {
	// 创建一个自定义路由
	mux := http.NewServeMux()
	// 调用
	// 创建静态文件服务器
	fs := http.FileServer(http.Dir("D:/Program/Mycode/Now/Mygo/Project/main/SnippetBox/ui/static"))
	// 去除前缀后从文件服务器中查找文件并返回
	// 不想让用户直接访问根目录可以检测访问路径并直接返回一个静态页面
	mux.Handle("/static/", http.StripPrefix("/static", fs))

	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/snippet/view", app.snippetView)
	mux.HandleFunc("/snippet/create", app.snippetCreate)

	// 使用中间件将当前mux下的所有路由都包装起来
	// 相当于是"重写"的在结构体中的方法
	// 最外层的中间件会第一个进行应用 类似于栈 first in first out
	return app.logRequest(secureHeaders(mux))
}
