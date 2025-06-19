package main

import (
	"gos/api"
	"net/http"
	"strconv"
	"strings"
)

func main() {
	http.HandleFunc("/api/chat", func(wr http.ResponseWriter, req *http.Request) {
		msg := req.URL.Query().Get("msg")

		if strings.HasPrefix(msg, "/prompt=") {
			api.Prompt = strings.TrimPrefix(msg, "/prompt=")
			wr.Write([]byte("系统提示词已设置为：" + api.Prompt))
			return
		} else if msg == "/prompt" {
			wr.Write([]byte("目前提示词: " + api.Prompt))
			return
		}

		if strings.HasPrefix(msg, "/token=") {
			ff := strings.TrimPrefix(msg, "/token=")
			switch ff {
			case "false":
				api.TokenSpent = false
				wr.Write([]byte("关闭Token回显"))
			case "true":
				api.TokenSpent = true
				wr.Write([]byte("开启Token回显"))
			}
			return
		} else if msg == "/token" {
			ff := api.TokenSpent
			wr.Write([]byte(strconv.FormatBool(ff)))
			return
		}

		if strings.HasPrefix(msg, "/weather=") {
			city := strings.TrimPrefix(msg, "/weather=")
			result := api.Weather(city)
			wr.Write([]byte(result))
			return
		}

		if msg == "/help" {
			wr.Write([]byte("/help"))
			return
		}
		if msg == "" {
			wr.WriteHeader(http.StatusBadRequest)
			wr.Write([]byte("缺少msg参数"))
			return
		}

		result := api.Chat(msg)
		wr.Write([]byte(result))
	})
	println(":5000")
	http.ListenAndServe(":5000", nil)
}
