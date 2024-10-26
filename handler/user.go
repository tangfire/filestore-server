package handler

import (
	dblayer "filestore-server/db"
	"filestore-server/util"
	"net/http"
	"os"
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

	if len(username) < 0 || len(password) < 5 {
		w.Write([]byte("Invild parameter"))
		return
	}

	enc_passwd := util.Sha1([]byte(password + pwd_salt))
	suc := dblayer.UserSignup(username, enc_passwd)
	if suc {
		w.Write([]byte("Signup success"))
	} else {
		w.Write([]byte("Signup fail"))
	}

}
