package main

import (
	"gos/api"
	"net/http"
)

func main() {
	http.HandleFunc("/api/chat", func(wr http.ResponseWriter, req *http.Request) {
		msg := req.URL.Query().Get("msg")
		if msg == "" {
			wr.WriteHeader(http.StatusBadRequest)
			wr.Write([]byte("缺少msg参数"))
			return
		}
		result := api.Chat(msg)
		wr.Write([]byte(result))
	})
	http.ListenAndServe(":5000", nil)
}
