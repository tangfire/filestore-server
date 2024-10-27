package handler

import (
	dblayer "filestore-server/db"
	"filestore-server/util"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
)

const (
	pwd_salt = "*#890"
)

// SignupHandler:处理用户注册请求
// 1.http GET请求，直接返回注册页面内容
// 2.校验参数的有效性
// 3.加密用户名密码
// 4.存入数据库tbl_user表及返回结果
func SignupHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		data, err := os.ReadFile("./static/view/signup.html")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Write(data)
		return
	}

	r.ParseForm()
	username := r.Form.Get("username")
	password := r.Form.Get("password")

	if len(username) < 3 || len(password) < 5 {
		w.Write([]byte("Invild parameter"))
		return
	}

	encPasswd := util.Sha1([]byte(password + pwd_salt))
	//fmt.Println("username = ", username, "encPasswd = ", encPasswd)
	suc := dblayer.UserSignup(username, encPasswd)

	if suc {
		w.Write([]byte("Signup success"))
	} else {
		w.Write([]byte("Signup fail"))
	}

}

// SignInHandler:登录接口
// 1.校验用户名及密码
// 2.生成访问凭证(token)
// 3.存储token到数据库tbl_user_token表
// 4.返回username,token,重定向url等信息
func SignInHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		// 返回html页面
		data, err := os.ReadFile("./static/view/signin.html")
		if err != nil {
			io.WriteString(w, "internet server error")
			return
		}
		io.WriteString(w, string(data))
	} else if r.Method == "POST" {
		r.ParseForm()
		username := r.Form.Get("username")
		password := r.Form.Get("password")
		//fmt.Println("SignIn received username: ", username) // 打印查看接收到的用户名
		encPasswd := util.Sha1([]byte(password + pwd_salt))
		// 1.校验用户名及密码
		pwdChecked := dblayer.UserSignin(username, encPasswd)
		if !pwdChecked {
			w.Write([]byte("FAILED"))
			return
		}
		// 2.生成访问凭证(token)
		token := GenToken(username)
		upRes := dblayer.UpdateToken(username, token)
		if !upRes {
			w.Write([]byte("FAILED"))
			return
		}

		// 3.登录成功后重定向到首页
		//w.Write([]byte("http://" + r.Host + "/static/view/home_01.html"))
		//http.Redirect(w, r, "/home", http.StatusFound)

		resp := util.RespMsg{
			Code: 0,
			Msg:  "OK",
			Data: struct {
				Location string
				Username string
				Token    string
			}{
				Location: "http://" + r.Host + "/home",
				Username: username,
				Token:    token,
			},
		}

		//fmt.Println("SignInHandler is finished")
		//w.Header().Set("Content-Type", "application/json")
		w.Write(resp.JSONBytes())
	}
}

// UserInfoHandler:查询用户信息
// 1.解析请求参数
// 2.验证token是否有效
// 3.查询用户信息
// 4.组装并且响应用户数据
func UserInfoHandler(w http.ResponseWriter, r *http.Request) {
	// 1.解释请求参数
	r.ParseForm()
	username := r.Form.Get("username")
	fmt.Println("UserInfoHandler received username = ", username) // 确认接收到的用户名
	//token := r.Form.Get("token")
	//// 2.验证token是否有效
	//isValidToken := IsTokenValid(username, token)
	//fmt.Println("isValidToken: ", isValidToken)
	//if !isValidToken {
	//	fmt.Println("token无效")
	//	w.WriteHeader(http.StatusForbidden)
	//	return
	//}
	// 3.查询用户信息
	fmt.Println("username = ", username)
	user, err := dblayer.GetUserInfo(username)
	if err != nil {
		fmt.Println("用户信息查询失败")
		w.WriteHeader(http.StatusForbidden)
		return
	}

	fmt.Println("user:", user.Username)
	fmt.Println("time", user.SignupAt)

	// 4.组装并且响应用户数据
	resp := util.RespMsg{
		Code: 0,
		Msg:  "OK",
		Data: user,
	}
	w.Write(resp.JSONBytes())
}

func GenToken(username string) string {
	// 40位字符:md5(username+timestamp+token_salt)+timestamp[:8]
	ts := fmt.Sprintf("%x", time.Now().Unix())
	tokenPrefix := util.MD5([]byte(username + ts + "_tokensalt"))
	return tokenPrefix + ts[:8]
}

//// IsTokenValid:token是否有效
//func IsTokenValid(token string) bool {
//	// TODO:判断token的时效性，是否过期
//
//	// TODO:从数据库表tbl_user_token查询username对应的user_token信息
//
//	// TODO:对比两个token是否一致
//
//	return true
//}

// IsTokenValid: token是否有效
func IsTokenValid(username, token string) bool {
	// 校验token的长度
	if len(token) != 40 {
		return false
	}

	// 获取用户的存储token
	storedToken, err := dblayer.GetUserToken(username)
	if err != nil {
		fmt.Println("Error retrieving token:", err)
		return false
	}

	// 比较传入的token与存储的token
	if token != storedToken {
		return false // token不一致
	}

	// 判断token的时效性，假设token的过期时间为8小时
	// 从token中提取时间戳
	ts := token[len(token)-8:] // 获取最后8个字符作为时间戳
	timestamp, err := strconv.ParseInt(ts, 16, 64)
	if err != nil {
		fmt.Println("Error parsing timestamp:", err)
		return false
	}

	// 判断当前时间与生成时间的差值
	if time.Now().Unix()-timestamp > 8*3600 {
		return false // token已过期
	}

	return true // token有效
}
