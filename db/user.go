package db

import (
	"database/sql"
	mydb "filestore-server/db/mysql"
	"fmt"
)

// UserSignup:通过用户名及密码完成user表的注册操作
func UserSignup(username string, password string) bool {

	stmt, err := mydb.DBConn().Prepare("insert into tbl_user (`user_name`,`user_pwd`) values (?,?)")
	if err != nil {
		fmt.Println("Failed to insert into tbl_user,err:" + err.Error())
		return false
	}

	defer stmt.Close()

	fmt.Println("UserSignup username:"+username, "password = "+password)
	ret, err := stmt.Exec(username, password)
	if err != nil {
		fmt.Println("Failed to insert into tbl_user,err:" + err.Error())
		return false
	}

	rowsAffected, err := ret.RowsAffected()
	if err != nil {
		fmt.Println("Failed to get rows affected, err:" + err.Error())
		return false
	}
	//fmt.Println("Rows affected:", rowsAffected)
	if rowsAffected > 0 {
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

type User struct {
	Username     string // 用户名
	Email        string
	Phone        string
	SignupAt     string // 注册时间
	LastActiveAt string
	Status       int
}

func GetUserInfo(username string) (User, error) {
	user := User{}

	stmt, err := mydb.DBConn().Prepare("select user_name,signup_at from tbl_user where user_name = ? limit 1")

	if err != nil {
		fmt.Println(err.Error())
		return user, err
	}

	defer stmt.Close()

	// 执行查询的操作
	err = stmt.QueryRow(username).Scan(&user.Username, &user.SignupAt)
	if err != nil {
		return user, err
	}

	return user, nil

}

// GetUserToken: 从数据库获取用户对应的 token
func GetUserToken(username string) (string, error) {
	//var db *sql.DB // 在实际应用中，需初始化数据库连接
	db := mydb.DBConn()
	var token string
	query := "SELECT user_token FROM tbl_user_token WHERE user_name = ?"
	err := db.QueryRow(query, username).Scan(&token)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("no token found for user: %s", username)
		}
		return "", err // 返回查询错误
	}

	return token, nil // 返回用户的 token
}
