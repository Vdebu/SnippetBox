package models

import (
	"database/sql"
	"errors"
	"time"
)

// 定义结构体存储数据库中提取出来的信息
type Snippet struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}

// SnippetModelInterface 定义接口用于解决模拟依赖注入时编译报错的问题
type SnippetModelInterface interface {
	Insert(title string, content string, expires int) (int, error)
	Get(id int) (*Snippet, error)
	Latest() ([]*Snippet, error)
}

// 定义模型存储数据库链接,有点类似依赖注入,在这之后定义方法
type SnippetModel struct {
	DB *sql.DB
}

// 都是直接返回错误由调用该函数的线程来处理错误 而不是在函数中直接处理错误

// 插入新的snippet
func (m *SnippetModel) Insert(title string, content string, expires int) (int, error) {
	// 使用占位符代替实际数据值
	//goland:noinspection SqlNoDataSourceInspection
	stmt := `INSERT INTO snippets(title,content,created,expires)
	VALUES(?,?,UTC_TIMESTAMP(),DATE_ADD(UTC_TIMESTAMP,INTERVAL ? DAY))`
	// 使用DB.Exec()执行SQL语句
	res, err := m.DB.Exec(stmt, title, content, expires)
	if err != nil {
		return 0, err
	}
	// 使用LastInsertId方法获取最后一次插入的id
	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	// id是int64 转换成int类型后正确返回
	return int(id), nil
}

// 输入id查询指定的snippet
//
//goland:noinspection SqlNoDataSourceInspection
func (m *SnippetModel) Get(id int) (*Snippet, error) {
	// 创建查询表达式
	stmt := `SELECT id,title,content,created,expires FROM snippets
	WHERE expires > UTC_TIMESTAMP AND id = ?`
	// 根据id获取当行的数据
	// 传入语句与用户输入的数据 返回一个指向sql.row的指针 存储从数据库查询得到的数据
	// 查询的语句与结果的提取是可以写到一起去的
	row := m.DB.QueryRow(stmt, id)
	// 使用数据结构尝试解析得到的数据
	s := &Snippet{}
	// 使用row.Scan()将查询到的数据写入到准备好的结构体里
	// 需要的参数是地址 要写入的参数个数必须与查询到的参数个数一致
	err := row.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
	if err != nil {
		// 先判断查询结果是不是空的 如果查询结果是空的会返回sql.ErrNoRows
		// 使用errors.Is()进行判断 会自动解包判断err
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		}
		return nil, err
	}
	// 将查找到的数据返回
	return s, nil
}

// 返回最新的十条snippet
//
//goland:noinspection SqlNoDataSourceInspection
func (m *SnippetModel) Latest() ([]*Snippet, error) {
	// 通过查询语句限制返回的内容数目
	stmt := `SELECT id,title,content,created,expires FROM snippets
	WHERE expires > UTC_TIMESTAMP()
	ORDER BY id DESC
	LIMIT 10`
	// 执行查询语句
	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	// 关闭数据流
	defer rows.Close()
	// 定义切片用于存储数据
	snippets := []*Snippet{}
	for rows.Next() {
		s := &Snippet{}
		// 尝试提取数据
		err = rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
		if err != nil {
			return nil, err
		}
		// 提取成功更新切片
		snippets = append(snippets, s)
	}
	// 在数据提取结束后判断过程中是否出错
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return snippets, nil
}
