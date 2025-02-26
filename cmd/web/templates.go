package main

import (
	"html/template"
	"io/fs"
	"path/filepath"
	"time"

	"SnippetBox.mikudayo.net/internal/models"
	"SnippetBox.mikudayo.net/ui"
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
	// 实现三方包中防止CSRF攻击的逻辑
	CSRFToken string
}

// 自定义时间格式化函数
func hunmanDate(t time.Time) string {
	// 如果传入的时间是空值直接返回空字符串
	if t.IsZero() {
		return ""
	}
	// 时间的初始化必须根据go的参考时间 不能随意更改
	// 在对时间格式化之前转换成utc时间
	return t.UTC().Format("2006-01-02 15:04:05")
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
	// pages, err := filepath.Glob("D:/Program/Mycode/Now/Mygo/Project/main/SnippetBox/ui/html/pages/*.html")

	// 使用fs.Glob从go embed的文件中提取出来文件名匹配的文件slice
	// 传入文件系统实例和一个 glob 模式，返回匹配的相对路径列表
	pages, err := fs.Glob(ui.Files, "html/pages/*.html")
	if err != nil {
		return nil, err
	}
	// 取出获取到的每一个文件路径
	for _, page := range pages {
		// 取出文件的名字用作key
		// home.tmpl.html
		name := filepath.Base(page)
		// 将需要渲染的模板打包成slice
		patterns := []string{
			"html/base.tmpl.html",
			"html/partials/*.html",
			page,
		}
		// 使用ParseFS()代替ParseFiles()从嵌入文件中对模板文件进行解析
		ts, err := template.New(name).Funcs(functions).ParseFS(ui.Files, patterns...)
		if err != nil {
			return nil, err
		}
		cache[name] = ts
	}
	// 返回编译好的所有内容
	return cache, nil
}
