package main

import (
	"database/sql"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"

	"SnippetBox.mikudayo.net/internal/models"

	_ "github.com/go-sql-driver/mysql"
)

type Application struct {
	// 向结构体注入自定义依赖 将处理器定义为结构体的方法
	errlog  *log.Logger
	infolog *log.Logger
	// 导入模型依赖
	snippets      *models.SnippetModel
	TemplateCache map[string]*template.Template
}

func main() {
	// 使用命令行参数来设置服务器端口 可以将设置信息存入环境变量 再调用命令行参数进行获取 最后传入预先定义好的变量中(结构体)
	addr := flag.String("addr", ":3939", "HTTP network address")
	// 严格区分大小写
	dsn := flag.String("dsn", "web:pass@/snippetbox?parseTime=true", "MySql data source name")
	// -addr=:4000指定参数 -help查看当前程序所有的可用参数

	// 使用前解析参数
	flag.Parse()

	// 自定义清晰的日志输出
	// 标准信息输出流 日期与时间
	infolog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	// 标准错误输出流 日期与时间与相关文件信息
	errlog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := openDB(*dsn)
	if err != nil {
		errlog.Fatal(err)
	}
	// 即使有时候程序会直接退出使defer的代码无法生效 添加关闭代码也是一个好的习惯
	defer db.Close()
	// 向结构体注入自定义依赖
	cache, err := newTemplateCache()
	if err != nil {
		errlog.Println(err)
		return
	}
	app := &Application{
		errlog:        errlog,
		infolog:       infolog,
		snippets:      &models.SnippetModel{DB: db},
		TemplateCache: cache,
	}
	// 自定义server结构体应用自定义的errlog否则在默认http遇到错误时还会调用原始的错误输出
	srv := &http.Server{
		// 改变默认端口
		Addr: *addr,
		// 改变输出错误日志的方法
		ErrorLog: errlog,
		// 设置自定义结构体中的处理器
		Handler: app.routes(),
	}
	infolog.Println("server start at", *addr, "...")
	// 设置了默认值之后使用新结构体的方法直接启动服务器
	err = srv.ListenAndServe()
	// 检查服务室是否会启动错误
	if err != nil {
		errlog.Fatal(err)
	}
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	// 使用Ping测试是否确实链接成功
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
