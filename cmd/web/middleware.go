package main

import "net/http"

// 添加保护网站安全的表头
func secureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 直接返回一个转化过的处理器
		w.Header().Set("Content-Security-Policy", "default-src 'self'; style-src 'self'fonts.googleapis.com; font-src fonts.gstatic.com")

		w.Header().Set("Referrer-Policy", "origin-when-cross-origin")
		w.Header().Set("X-Content-Type-0ptions", "nosniff")
		w.Header().Set("X-Frame-0ptions", "deny")
		w.Header().Set("X-XSS-Protection", "0")
		// 调用原始的处理器
		next.ServeHTTP(w, r)
	})
}
