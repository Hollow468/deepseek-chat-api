package api

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

var ImgModel = "Kwai-Kolors/Kolors"

func CreateImg(prompt string, negativePrompt ...string) string {
	url := "https://api.siliconflow.cn/v1/images/generations"

	baseJSON := `{
  "model": "` + ImgModel + `",
  "prompt": "` + prompt + `",
  "image_size": "1024x1024",
  "batch_size": 1,
  "num_inference_steps": 20,
  "guidance_scale": 7.5`

	if len(negativePrompt) > 0 && negativePrompt[0] != "" {
		baseJSON += `,
  "negative_prompt": "` + negativePrompt[0] + `"`
	}

	baseJSON += `
}`

	jstr := baseJSON

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
