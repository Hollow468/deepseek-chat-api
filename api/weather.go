package api

import (
	"log"
	"net/http"
	"io"
	"encoding/json"
)

type WeatherResponse struct {
	Status string `json:"status"`
	Lives []weatherInfo `json:"lives"`
}
type weatherInfo struct {
	City string `json:"city"`
	Weather string `json:"weather"`
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
	
	jsonBytes, err := json.Marshal(result)
	if err != nil {
		log.Fatal(err)
	}
	return string(jsonBytes)
}
