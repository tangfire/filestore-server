package mysql

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql" // 注意这个导入
	"log"
	"os"
)

// 通过sql.DB来管理数据库连接对象
var db *sql.DB

func init() {
	// 错误处理: var err error 在这里声明一个局部变量，然后使用 db, err = sql.Open(...) 来赋值。
	// 这样可以确保 db 变量是全局的，而 err 是局部的。

	// 使用全局变量: 保证 db 是全局变量，在 DBConn() 函数中返回的就是正确初始化的 db。

	var err error // 声明一个局部变量 err

	//通过sql.Open来创建协程安全的sql.DB对象
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

// ParseRows: 解析 *sql.Rows 结果并返回切片，元素为 map[string]interface{}
// 解释：
//
//	1、获取列名：通过 rows.Columns() 获取列名，并存储在 columns 切片中。
//	2、准备存储数据：创建一个切片 values 用于存储每一行的数据，使用 new(interface{}) 以支持任意类型的值。
//	3、循环扫描行：使用 rows.Next() 遍历结果集的每一行，并用 rows.Scan() 将行数据存入 values。
//	4、构建映射：为每一行数据创建一个 map[string]interface{}，将列名作为键，列值作为值，最后将这个映射添加到结果切片中。
//	5、错误处理：检查 rows.Err() 确保在遍历过程中没有发生错误。
//
// 这样，你就能通过 mydb.ParseRows(rows) 将 *sql.Rows 转换为易于处理的数据结构。
func ParseRows(rows *sql.Rows) []map[string]interface{} {
	columns, err := rows.Columns()
	if err != nil {
		log.Println("Failed to get columns:", err)
		return nil
	}

	// 创建一个切片，用于存储每一行的结果
	var result []map[string]interface{}

	// 创建一个切片用于存储行数据
	values := make([]interface{}, len(columns))
	for i := range values {
		values[i] = new(interface{})
	}

	for rows.Next() {
		// 扫描行数据到 values 切片
		if err := rows.Scan(values...); err != nil {
			log.Println("Failed to scan row:", err)
			continue
		}

		// 创建一个 map 存储当前行的数据
		rowData := make(map[string]interface{})
		for i, col := range columns {
			rowData[col] = *(values[i].(*interface{}))
		}

		result = append(result, rowData)
	}

	if err := rows.Err(); err != nil {
		log.Println("Error occurred during rows iteration:", err)
	}

	return result
}
