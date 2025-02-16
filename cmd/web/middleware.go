package main

import (
	"fmt"
	"net/http"
)

// 添加保护网站安全的表头
func secureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 直接返回一个转化过的处理器
		// 响应头写错了会导致页面无法正常加载
		w.Header().Set("Content-Security-Policy", "default-src 'self'; style-src 'self' fonts.googleapis.com; font-src fonts.gstatic.com")

		w.Header().Set("Referrer-Policy", "origin-when-cross-origin")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "deny")
		w.Header().Set("X-XSS-Protection", "0")
		// 调用原始的处理器
		next.ServeHTTP(w, r)
	})
}

// 记录下每一个请求的信息
func (app *Application) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 输出请求的具体信息
		app.infolog.Printf("%s - %s %s %s", r.RemoteAddr, r.Proto, r.Method, r.URL.RequestURI())

		next.ServeHTTP(w, r)
	})
}

// 检查是否发生过panic输出人性化的提示
func (app *Application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 创建一个deferred函数确保在请求后必定会检查是否有错误发生并执行相关逻辑
		defer func() {
			// 使用内置函数判断是否有panic发生
			if err := recover(); err != nil {
				// 有错误发生使用humanOutPut
				// 向响应体写入链接关闭的消息
				w.Header().Set("Connection", "close")
				// 将遇到的错误包装返回
				app.serverError(w, fmt.Errorf("%s", err))
			}
		}()
		next.ServeHTTP(w, r)
	})
}
