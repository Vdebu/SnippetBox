package main

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/go-playground/form/v4"
)

// 将函数输出错误信息的权限大部分移交给helper(app.errlog,)

// 输出错误信息与栈追踪(在那个goroutine中调用的这个函数)
func (app *Application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	// 在日志中输出当前调用函数的goroutine
	// ERROR   2025/02/12 13:56:08 helpers.go:13:
	// ->ERROR   2025/02/12 13:56:19 handlers.go:36
	app.errlog.Output(2, trace)

	// 输出内部服务器错误的信息 statusText会根据传入的代码自动生成可读的错误信息 s
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

// 输出客户端(http)的错误信息,一般是由用户自己造成的
func (app *Application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

// 返回404notfound 通过包装clientError实现
func (app *Application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}

// 用于渲染各个网页
func (app *Application) render(w http.ResponseWriter, status int, page string, data *TemplateData) {
	// 从模板缓存中获取当前请求页面的模板
	ts, ok := app.templateCache[page]
	if !ok {
		// 若当前请求的页面模板不存在
		// 定义一个新的错误并报告
		err := fmt.Errorf("the template %s does not exist", page)
		app.serverError(w, err)
		return
	}
	// 创建一个字节类型的缓冲
	buf := new(bytes.Buffer)
	// 将获取到的模板写入缓冲查看是否成功
	err := ts.ExecuteTemplate(buf, "base", data)
	if err != nil {
		app.serverError(w, err)
		return
	}
	// 向缓冲写入成功没有报错

	// 向html头部写入提供的状态码
	w.WriteHeader(status)
	// 向响应体写入数据
	buf.WriteTo(w)
}

// 通过解码将读取到的内容写入传入的结构体(指针)中
func (app *Application) decodePostForm(r *http.Request, dst any) error {
	// 先对request进行解析
	err := r.ParseForm()
	if err != nil {
		return err
	}
	// 尝试将解析得到的数据按结构体中的要求解码
	err = app.formDecoder.Decode(dst, r.PostForm)
	if err != nil {
		// 只处理传入的数据存储器不是指针的错误
		// 定义出decode会返回的无效的解码错误用Errors的方法进行解码
		var invalidDecoderError *form.InvalidDecoderError
		if errors.As(err, &invalidDecoderError) {
			// 直接panic交给recover中间件进行处理
			panic(err)
		}
		// 如果并非无效值错误就直接返回
		return err
	}
	return nil
}

// 默认出事TemplateData结构体中的时间字段(每个网页底部都需要显示时间)
func (app *Application) newTemplateData(r *http.Request) *TemplateData {
	return &TemplateData{
		// 获取当前的年份
		CurrentYear: time.Now().Year(),
	}
}
