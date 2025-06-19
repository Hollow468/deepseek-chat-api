package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

var Model = "deepseek-ai/DeepSeek-V3"
var TokenSpent = false

func Chat(message string) string {
	url := "https://api.siliconflow.cn/v1/chat/completions"
	jsonStr := `{
  "model": "` + Model + `",
  "messages": [
    {
      "role": "user",
      "content": "` + message + `"
    }
  ],
  "stream": false,
  "max_tokens": 512,
  "enable_thinking": false,
  "thinking_budget": 4096,
  "min_p": 0.05,
  "stop": null,
  "temperature": 0.7,
  "top_p": 0.7,
  "top_k": 50,
  "frequency_penalty": 0.5,
  "n": 1,
  "response_format": {
    "type": "text"
  },
  "tools": [
    {
      "type": "function",
      "function": {
        "description": "<string>",
        "name": "<string>",
        "parameters": {},
        "strict": false
      }
    }
  ]
}`

	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(jsonStr)))
	if err != nil {
		panic(err)
	}
	req.Header.Set("Authorization", "Bearer "+SiliflowKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	var result struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}
	var usage struct {
		Usage struct {
			TotalTokens int `json:"total_tokens"`
		} `json:"usage"`
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(body, &usage)
	if err != nil {
		panic(err)
	}
	if len(result.Choices) > 0 {
		if TokenSpent {
			return result.Choices[0].Message.Content + "\ntotal_tokens: " + fmt.Sprint(usage.Usage.TotalTokens)
		} else {
			return result.Choices[0].Message.Content
		}
	}
	return ""
}
