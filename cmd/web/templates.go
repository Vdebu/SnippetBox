package main

import (
	"html/template"
	"path/filepath"
	"time"

	"SnippetBox.mikudayo.net/internal/models"
)

// 用于存储渲染网页所需要用到的数据
type TemplateData struct {
	// 显示在网页下方的年份
	CurrentYear int
	// 渲染的时候需要将数据注入 所以将snippet都嵌入到TemplateData中 主依赖只包含数据库与操作方法
	// 渲染网页要用到的主体
	Snippet  *models.Snippet
	Snippets []*models.Snippet
	// 用于存储用户输入的错误信息重新渲染页面
	Form  any
	Flash string
	// 存储当前用户是否登入的信息
	IsAuthenticated bool
}

// 自定义时间格式化函数
func hunmanDate(t time.Time) string {
	// 时间的初始化必须根据go的参考时间 不能随意更改
	return t.Format("2006-01-02 15:04:05")
}

// 创建template.FuncMap用于存储自定义函数
var functions = template.FuncMap{
	"humanDate": hunmanDate,
}

// 将网页模板渲染并存储到内存中 提高运行效率
func newTemplateCache() (map[string]*template.Template, error) {
	// 初始化容器用于存储
	cache := map[string]*template.Template{}
	// 尝试获取指定目录下的同类型文件
	pages, err := filepath.Glob("D:/Program/Mycode/Now/Mygo/Project/main/SnippetBox/ui/html/pages/*.html")
	if err != nil {
		return nil, err
	}
	// 取出获取到的每一个文件路径
	for _, page := range pages {
		// 取出文件的名字用作key
		// home.tmpl.html
		name := filepath.Base(page)
		// 先解析基础的模板
		// 为了应用函数到模板里 需要使用new方法创建一个新的template.Template对象
		// 函数的应用必须与模板解析之前
		ts, err := template.New(name).Funcs(functions).ParseFiles("D:/Program/Mycode/Now/Mygo/Project/main/SnippetBox/ui/html/base.tmpl.html")
		if err != nil {
			return nil, err
		}
		// 调用ts.ParseGlob()解析其他的基板(nav.tmpl.html)
		ts, err = ts.ParseGlob("D:/Program/Mycode/Now/Mygo/Project/main/SnippetBox/ui/html/partials/*.tmpl.html")
		if err != nil {
			return nil, err
		}
		// 处理完所有的基板开始渲染页面
		ts, err = ts.ParseFiles(page)
		if err != nil {
			return nil, err
		}
		// 预编译成功将模板存入字典中
		cache[name] = ts
	}
	// 返回编译好的所有内容
	return cache, nil
}
