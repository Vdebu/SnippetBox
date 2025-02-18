package main

// 定义所有的处理器

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"SnippetBox.mikudayo.net/internal/models"
	"github.com/julienschmidt/httprouter"
)

func (app *Application) home(w http.ResponseWriter, r *http.Request) {
	// 指定"/"的逻辑 防止预料之外的访问
	if r.URL.Path != "/" {
		// 直接使用not found方法写入信息
		app.notFound(w)
		return
	}

	// 每一个请求都是一个独立的goroutine 这里的panic不会导致server崩溃
	// panic("something go wrong...")

	snippets, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, err)
	}
	// 先初始化默认数据再初始化查询得到的数据
	data := app.newTemplateData(r)
	data.Snippets = snippets
	// 使用render()进行home.tmpl.html渲染
	app.render(w, http.StatusOK, "home.tmpl.html", data)
	// w.Write([]byte("mikudayoooo"))
}

func (app *Application) snippetView(w http.ResponseWriter, r *http.Request) {
	// 获取url ?后的id用于数据库查询
	// 抽取为string并用于转换判断是否为正值
	// http://localhost:3939/snippet/view?id=39
	// id, err := strconv.Atoi(r.URL.Query().Get("id"))
	// w.Write([]byte("Display a specific miku..."))
	// 使用新的方法获取url中的值
	params := httprouter.ParamsFromContext(r.Context())
	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil || id < 0 {
		// http.NotFound(w, r)
		app.notFound(w)
		return
	}
	snippet, err := app.snippets.Get(id)
	if err != nil {
		// 判断是否是ErrNoRecord 这里要通过包名调用自己定义的错误model.ErrNoRecord
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
			return
		}
		app.serverError(w, err)
		return
	}
	data := app.newTemplateData(r)
	data.Snippet = snippet
	app.render(w, http.StatusOK, "view.tmpl.html", data)
	// 将搜索到的内容直接输出到响应体
	// fmt.Fprintf(w, "Display a specific miku %v...", snippet)
}

func (app *Application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	// 指定request的方法必须是POST
	// if r.Method != http.MethodPost {
	// 	// 使用定义在net/http中的常量来代替直接的字符串

	// 	// 添加新的header展示更详细的信息
	// 	w.Header().Set("Allow", http.MethodPost)

	// 	// 写入状态码405 表示request方法不被允许
	// 	// w.WriteHeader(http.StatusMethodNotAllowed)
	// 	// w.Write([]byte("method not allowed..."))

	// 	// 一般不会直接调用writeheader与write而是通过别的函数间接调用
	// 	// http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	// 	app.clientError(w, http.StatusMethodNotAllowed)
	// 	return
	// }

	// 使用三方库可以在初始化处理器的时候直接制定请求方法
	// 后续用于返回一个填写信息的网页

	// 创建数据结构体用于后续写入 因为这里需要数据中的时间(只初始化了时间)所以才需要再初始化这个结构体
	// 这也保证了代码的结构统一
	data := app.newTemplateData(r)

	app.render(w, http.StatusOK, "create.tmpl.html", data)
	w.Write([]byte("Createa a new miku..."))

}

func (app *Application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {
	// 只接受Post请求的处理器
	title := "mikudayo"
	content := "mikudayo\nmikudayo\nmikudayo"
	expires := 7
	id, err := app.snippets.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, err)
		return
	}
	// 创建成功后将用户重定向到最新创建的snippet
	// curl -iL -X POST http://localhost:3939/snippet/create
	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}
