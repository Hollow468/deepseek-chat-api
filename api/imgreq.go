package api

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

func CreateImg(prompt string) string {
	url := "https://api.siliconflow.cn/v1/images/generations"
	jstr := `{
  "model": "Kwai-Kolors/Kolors",
  "prompt": "` + prompt + `",
  "image_size": "1024x1024",
  "batch_size": 1,
  "num_inference_steps": 20,
  "guidance_scale": 7.5
}`

	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(jstr)))
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

	type a struct {
		A []struct {
			Url string `json:"url"`
		} `json:"images"`
	}

	var res a
	err = json.Unmarshal(body, &res)
	if err != nil {
		panic(err)
	}

	if len(res.A) > 0 {
		return res.A[0].Url
	}
	return ""
}
