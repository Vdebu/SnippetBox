package main

import (
	"html/template"
	"path/filepath"

	"SnippetBox.mikudayo.net/internal/models"
)

type TemplateData struct {
	Snippet  *models.Snippet
	Snippets []*models.Snippet
}

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
		// 将路径导入到预先准备好的模板包中
		// files := []string{
		// 	"D:/Program/Mycode/Now/Mygo/Project/main/SnippetBox/ui/html/base.tmpl.html",
		// 	"D:/Program/Mycode/Now/Mygo/Project/main/SnippetBox/ui/html/partials/nav.tmpl.html",
		// 	page,
		// }
		// 先解析基础的模板
		ts, err := template.ParseFiles("D:/Program/Mycode/Now/Mygo/Project/main/SnippetBox/ui/html/base.tmpl.html")
		if err != nil {
			return nil, err
		}
		// 调用ts.ParseGlob()解析其他的基板
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
