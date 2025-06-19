package api

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

func ModelCheck(model string) bool {
	url := "https://api.siliconflow.cn/v1/models"

	req, err := http.NewRequest("GET", url, nil)
	req.Header.Add("Authorization", "Bearer "+SiliflowKey)

	if err != nil {
		log.Fatal(err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	type Model struct {
		ID string `json:"id"`
	}
	type ModelList struct {
		Data []Model `json:"data"`
	}

	var models ModelList
	err = json.Unmarshal(body, &models)
	if err != nil {
		log.Fatal(err)
	}

	for _, m := range models.Data {
		if m.ID == model {
			return true
		}
	}
	return false
}
