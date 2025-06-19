package api

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
)

type balanceData struct {
	TotalBalance string `json:"totalBalance"`
}

type balanceOuter struct {
	Data balanceData `json:"data"`
}

func Balance() string {
	req, err := http.NewRequest("GET", "https://api.siliconflow.cn/v1/user/info", nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Add("Authorization", "Bearer "+SiliflowKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var result balanceOuter
	err = json.Unmarshal(body, &result)
	if err != nil {
		log.Fatal(err)
	}
	balance, err := strconv.ParseFloat(result.Data.TotalBalance, 64)
	if err != nil {
		log.Fatal(err)
	}
	return fmt.Sprintf("当前余额为：%f", balance)
}
