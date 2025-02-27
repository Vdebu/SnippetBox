package main

import (
	"SnippetBox.mikudayo.net/internal/assert"
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSecureHeader(t *testing.T) {
	// 创建新的记录器记录待测试的中间件的响应结果
	rr := httptest.NewRecorder()

	r, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	// 创建一个处理器 后续用中间件将其进行包装
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})
	// 使用中间件对简单处理器进行包装并执行ServeHTTP进行测试
	// 传入用于测试定义的响应体与请求
	secureHeaders(next).ServeHTTP(rr, r)
	rs := rr.Result()
	// 使用键值提取特定的表头进行比较查看是否写入成功
	expectedValue := "default-src 'self'; style-src 'self' fonts.googleapis.com; font-src fonts.gstatic.com"
	assert.Equal(t, rs.Header.Get("Content-Security-Policy"), expectedValue)
	expectedValue = "origin-when-cross-origin"
	assert.Equal(t, rs.Header.Get("Referrer-Policy"), expectedValue)
	expectedValue = "nosniff"
	assert.Equal(t, rs.Header.Get("X-Content-Type-Options"), expectedValue)
	expectedValue = "deny"
	assert.Equal(t, rs.Header.Get("X-Frame-Options"), expectedValue)
	expectedValue = "0"
	assert.Equal(t, rs.Header.Get("X-XSS-Protection"), expectedValue)

	// 检查状态码是否符合预期
	assert.Equal(t, rs.StatusCode, http.StatusOK)

	// 检查响应体中的内容
	defer rs.Body.Close()
	body, err := io.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}
	// 去除响应体中的空格
	bytes.TrimSpace(body)
	// 判断是否一致
	assert.Equal(t, string(body), "OK")
}
