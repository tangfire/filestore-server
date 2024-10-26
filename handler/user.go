package handler

import (
	dblayer "filestore-server/db"
	"filestore-server/util"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

const (
	pwd_salt = "*#890"
)

// SignupHandler:处理用户注册请求
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
	suc := dblayer.UserSignup(username, encPasswd)
	if suc {
		w.Write([]byte("Signup success"))
	} else {
		w.Write([]byte("Signup fail"))
	}

}

// SignInHandler:登录接口
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
		//w.Write([]byte("http://" + r.Host + "/static/view/home.html"))
		http.Redirect(w, r, "/home", http.StatusFound)
	}

}

func GenToken(username string) string {
	// 40位字符:md5(username+timestamp+token_salt)+timestamp[:8]
	ts := fmt.Sprintf("%x", time.Now().Unix())
	tokenPrefix := util.MD5([]byte(username + ts + "_tokensalt"))
	return tokenPrefix + ts[:8]
}
