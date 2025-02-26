package ui

import "embed"

// 在程序编译的时候嵌入静态文件
// 语句必须写在一个全局变量的上边(只能向全局变量嵌入文件 会忽略.或/开头的文件)
// 可以用在FileServer里面

//go:embed "html" "static"
var Files embed.FS
