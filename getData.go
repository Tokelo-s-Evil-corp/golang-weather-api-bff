package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
)

type CurrentWeather struct {
	Time     string  `json:"time"`
	Interval int     `json:"interval"`
	Temp2M   float64 `json:"temperature_2m"`
}

type CurrentUnits struct {
	Time     string `json:"time"`
	Interval string `json:"interval"`
	Temp     string `json:"temperature_2m"`
}

type Item struct {
	CurrentWeather CurrentWeather `json:"current"`
	CurrentUnits   CurrentUnits   `json:"current_units"`
}

var myData CurrentWeather

func makeDataReq() {

	//make get request for data with parameters

	baseURL := "https://api.open-meteo.com/v1/forecast"

	params := url.Values{}
	params.Add("latitude", "-29.3167")
	params.Add("longitude", "27.4833")
	params.Add("current", "temperature_2m")
	params.Add("hourly", "temperature_2m")

	queryString := params.Encode()

	url := baseURL + "?" + queryString

	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		fmt.Print(err)
		return
	}

	//send request
	client := &http.Client{}

	resp, err := client.Do(req)

	if err != nil {
		fmt.Print(err)
		return
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		fmt.Print(err)
		return
	}

	var data Item

	error := json.Unmarshal(body, &data)

	if error != nil {
		fmt.Print(error)
	}

	myData = data.CurrentWeather

}

func endpoints() {
	router := gin.Default()

	router.GET("/currentweather", getWeather)

	router.Run("localhost:8000")
}

func getWeather(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, myData)
}
