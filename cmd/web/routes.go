package main

import "net/http"

func (app *Application) routes() *http.ServeMux {
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
	
	return mux
}
