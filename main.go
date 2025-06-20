package main

import (
	"flag"
	"gos/api"
	"gos/load"
	"net/http"
	"strconv"
	"strings"
)

func main() {
	port := flag.String("port", "5000", "http监听端口")
	flag.Parse()

	cfg, err := load.LoadConfig()
	if err == nil {
		api.Model = cfg.Model
		api.Prompt = cfg.Prompt
		api.TokenSpent = cfg.TokenSpent
		println("config loaded")
	}

	http.HandleFunc("/api/chat", func(wr http.ResponseWriter, req *http.Request) {
		msg := req.URL.Query().Get("msg")

		// Debug信息
		println("收到msg:", msg)

		if strings.HasPrefix(msg, "/prompt ") {
			api.Prompt = strings.TrimPrefix(msg, "/prompt ")
			wr.Write([]byte("系统提示词已设置为：" + api.Prompt))
			return
		} else if msg == "/prompt" {
			wr.Write([]byte("目前提示词: " + api.Prompt))
			return
		}

		if strings.HasPrefix(msg, "/token ") {
			ff := strings.TrimPrefix(msg, "/token ")
			if n, err := strconv.Atoi(ff); err == nil {
				api.MaxTokens = n
				wr.Write([]byte("max_tokens已设置为:" + ff))
				return
			}
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

		if strings.HasPrefix(msg, "/weather ") {
			city := strings.TrimPrefix(msg, "/weather ")
			result := api.Weather(city)
			wr.Write([]byte(result))
			return
		}

		if strings.HasPrefix(msg, "/balance") {
			wr.Write([]byte(api.Balance()))
			return
		}

		if strings.HasPrefix(msg, "/model set ") {
			m := strings.TrimPrefix(msg, "/model set ")
			if api.ModelCheck(m) {
				api.Model = m
				wr.Write([]byte("设置成功:" + api.Model))
				return
			} else {
				wr.Write([]byte("模型无效"))
				return
			}
		} else if msg == "/model" {
			wr.Write([]byte("当前模型:" + api.Model))
			return
		}

		if msg == "/model list" || msg == "/model ls" {
			models := api.GetAllModels()
			wr.Write([]byte("可用模型：\n" + strings.Join(models, "\n")))
			return
		}

		if strings.HasPrefix(msg, "/img ") {
			ct := strings.TrimPrefix(msg, "/img ")

			args := strings.Fields(ct)
			if len(args) == 0 {
				wr.Write([]byte("请提供图片描述"))
				return
			}

			var prompt string
			var negativePrompt string
			var otherParams []string

			for i := 0; i < len(args); i++ {
				arg := args[i]
				if arg == "--negative" {
					if i+1 < len(args) {
						negativeArgs := args[i+1:]
						negativePrompt = strings.Join(negativeArgs, " ")
						break
					}
				} else if strings.HasPrefix(arg, "--") {
					otherParams = append(otherParams, arg)
				} else {
					if prompt == "" {
						prompt = arg
					} else {
						prompt += " " + arg
					}
				}
			}

			if negativePrompt == "" {
				prompt = ct
			}

			//debug
			println("正面提示词:", prompt)
			if negativePrompt != "" {
				println("负面提示词:", negativePrompt)
			}
			if len(otherParams) > 0 {
				println("其他参数:", strings.Join(otherParams, " "))
			}

			rq := api.CreateImg(prompt, negativePrompt)
			wr.Write([]byte(rq))
			return
		}

		if strings.HasPrefix(msg, "/model img ") {
			model := strings.TrimPrefix(msg, "/model img ")
			api.ImgModel = model
			wr.Write([]byte("图片模型已设置为: " + api.ImgModel))
			return
		}

		if msg == "/help" {
			helpText := `可用命令：
/help - 显示帮助信息
/prompt xxx - 设置系统提示词
/prompt - 查看当前提示词
/token true|false - 开启/关闭Token回显
/token 数字 - 设置max_tokens数量
/token - 查看当前Token回显状态
/weather 城市名 - 查询天气
/balance - 查看账户余额
/model set 模型名 - 设置模型
/model - 查看当前模型
/model list|ls - 查看所有可用模型
/model img xx - 配置生图模型
/save - 保存当前配置到文件
/clear - 清空聊天历史
/config list - 查看所有配置项
/img 描述 - 生成图片
/img 描述 --negative 负面描述 - 生成图片（带负面提示词）`
			wr.Write([]byte(helpText))
			return
		}

		if msg == "/save" {
			cfg := load.ChatConfig{
				Model:      api.Model,
				Prompt:     api.Prompt,
				TokenSpent: api.TokenSpent,
			}
			err := load.SaveConfig(cfg)
			if err != nil {
				wr.Write([]byte("保存失败: " + err.Error()))
			} else {
				wr.Write([]byte("配置已保存"))
			}
			return
		}

		if msg == "/clear" {
			api.ClearHistory()
			wr.Write([]byte("历史记录已清空"))
			return
		}

		if msg == "" {
			wr.WriteHeader(http.StatusBadRequest)
			wr.Write([]byte("缺少msg参数"))
			return
		}

		result := api.Chat(msg)
		println("返回结果:", result)
		wr.Write([]byte(result))
	})

	println(":" + *port)
	http.ListenAndServe(":"+*port, nil)
}
