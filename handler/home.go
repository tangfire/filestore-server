package handler

import (
	"io"
	"net/http"
	"os"
)

func GoHomeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		data, err := os.ReadFile("./static/view/home.html")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		io.WriteString(w, string(data))
		return
	}
}
