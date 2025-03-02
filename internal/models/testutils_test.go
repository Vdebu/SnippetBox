package models

import (
	"database/sql"
	"os"
	"testing"
)

// 链接测试用的数据库
func newTestDB(t *testing.T) *sql.DB {
	dst := "testweb:pass@/test_snippetbox?parseTime=true&multiStatements=true"
	// 尝试链接数据库
	db, err := sql.Open("mysql", dst)
	if err != nil {
		t.Fatal(err)
	}
	// 读取建立数据库表格与插入数据的脚本
	script, err := os.ReadFile("./testdata/setup.sql")
	if err != nil {
		t.Fatal(err)
	}
	// 执行脚本中的sql语句
	_, err = db.Exec(string(script))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		script, err := os.ReadFile("./testdata/teardown.sql")
		if err != nil {
			t.Fatal(err)
		}
		// 执行脚本清理测试数据库中的数据
		_, err = db.Exec(string(script))
		if err != nil {
			t.Fatal(err)
		}
		// 关闭数据库的连接
		db.Close()
	})
	return db
}
