package api

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

type WeatherResponse struct {
	Lives []weatherInfo `json:"lives"`
}
type weatherInfo struct {
	City        string `json:"city"`
	Weather     string `json:"weather"`
	Temperature string `json:"temperature"`
}

func Weather(city string) string {
	url := "https://restapi.amap.com/v3/weather/weatherInfo?key=" + WeatherKey + "&city=" + city

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}

	response, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	content, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	var result WeatherResponse
	err = json.Unmarshal(content, &result)
	if err != nil {
		log.Fatal(err)
	}

	if len(result.Lives) == 0 || result.Lives[0].City == "" {
		return "未找到该城市天气信息"
	}
	info := result.Lives[0]
	return info.City + " 天气：" + info.Weather + "，温度：" + info.Temperature + "℃"
}
