package mysql

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql" // 注意这个导入
	"os"
)

var db *sql.DB

func init() {
	// 错误处理: var err error 在这里声明一个局部变量，然后使用 db, err = sql.Open(...) 来赋值。
	// 这样可以确保 db 变量是全局的，而 err 是局部的。

	// 使用全局变量: 保证 db 是全局变量，在 DBConn() 函数中返回的就是正确初始化的 db。

	var err error                                                                            // 声明一个局部变量 err
	db, err = sql.Open("mysql", "root:8888.216@tcp(127.0.0.1:3306)/fileserver?charset=utf8") // 不要在这里使用 := 赋值
	if err != nil {
		fmt.Println("Failed to open database, err:", err)
		os.Exit(1)
	}
	db.SetMaxOpenConns(1000)
	errPing := db.Ping()
	if errPing != nil {
		fmt.Println("Failed to connect to mysql,err:" + errPing.Error())
		os.Exit(1)
	}

}

// DBConn:返回数据库连接对象
func DBConn() *sql.DB {
	return db
}
