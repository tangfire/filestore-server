package handler

import (
	"fmt"
	"net/http"
)

// HTTPInterceptor:http请求拦截器
func HttpInterceptor(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()

		username := r.Form.Get("username")
		token := r.Form.Get("token")

		if len(username) < 0 || !IsTokenValid(username, token) {

			w.WriteHeader(http.StatusForbidden)
			return
		}
		fmt.Println("HttpInterceptor Acc")
		h(w, r)
	}
}
