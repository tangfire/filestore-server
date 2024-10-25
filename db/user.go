package db

import (
	mydb "filestore-server/db/mysql"
	"fmt"
)

// UserSignup:通过用户名及密码完成user表的注册操作
func UserSignup(username string, password string) bool {
	stmt, err := mydb.DBConn().Prepare("insert ignore into tbl_user (`user_name`,`user_pwd`) values (?,?)")
	if err != nil {
		fmt.Println("Failed to insert into tbl_user,err:" + err.Error())
		return false
	}

	defer stmt.Close()

	ret, err := stmt.Exec(username, password)
	if err != nil {
		fmt.Println("Failed to insert into tbl_user,err:" + err.Error())
		return false
	}

	if rowsAffected, err := ret.RowsAffected(); nil == err && rowsAffected > 0 {
		return true
	}

	return false
}

// UserSignin:判断密码是否一致
func UserSignin(username string, encpwd string) bool {
	stmt, err := mydb.DBConn().Prepare("select * from tbl_user where user_name=? limit 1")
	if err != nil {
		fmt.Println("Failed to select tbl_user,err:" + err.Error())
		return false
	}

	rows, err := stmt.Query(username)
	if err != nil {
		fmt.Println("Failed to select tbl_user,err:" + err.Error())
		return false
	} else if rows == nil {
		fmt.Println("username not found:" + username)
		return false
	}

	pRows := mydb.ParseRows(rows)
	if len(pRows) > 0 && string(pRows[0]["user_pwd"].([]byte)) == encpwd {
		return true
	}

	return false

}

// UpdateToken:刷新用户登录的token
func UpdateToken(username string, token string) bool {
	stmt, err := mydb.DBConn().Prepare("replace into tbl_user_token (`user_name`,`user_token`) values(?,?)")
	if err != nil {
		fmt.Println("Failed to replace tbl_user_token,err:" + err.Error())
		return false
	}

	defer stmt.Close()

	_, err = stmt.Exec(username, token)
	if err != nil {
		fmt.Println("Failed to replace tbl_user_token,err:" + err.Error())
		return false
	}

	return true

}
