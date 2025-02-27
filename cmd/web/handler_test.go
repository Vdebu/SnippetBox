package main

import (
	"SnippetBox.mikudayo.net/internal/assert"
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPing(t *testing.T) {
	// 初始化responserecoder用于后续传入待测试的处理器中记录结果
	rr := httptest.NewRecorder()
	// 初始化一个简单的http.request
	r, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	// 测试处理器
	ping(rr, r)
	// 获取响应体的内容
	rs := rr.Result()
	// 检查状态码是否与预期的一致
	assert.Equal(t, rs.StatusCode, http.StatusOK)

	// 检查响应体中写入的数据是否是OK
	defer rs.Body.Close()
	body, err := io.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}
	// 去除响应体中的空格用于与预先写入的数据进行比较
	bytes.TrimSpace(body)
	assert.Equal(t, string(body), "OK")
}
