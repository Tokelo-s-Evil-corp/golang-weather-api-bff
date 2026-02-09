package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type CurrentWeather struct {
	Time                      string  `json:"time"`
	Interval                  int     `json:"interval"`
	Temp2M                    float64 `json:"temperature_2m"`
	WindSpeed                 float64 `json:"wind_speed_10m"`
	Humidity                  int64   `json:"relative_humidity_2m"`
	Precipitation             float64 `json:"precipitation"`
	Precipitation_Probability int     `json:"precipitation_probability"`
	WeatherCode               int     `json:"weather_code"`
	WindDirection             int     `json:"wind_direction_10m"`
	ApparentTemperature       float64 `json:"apparent_temperature"`
	CloudCover                int     `json:"cloud_cover"`
}

type CurrentUnits struct {
	Time                      string `json:"time"`
	Interval                  string `json:"interval"`
	Temp                      string `json:"temperature_2m"`
	WindSpeed                 string `json:"wind_speed_10m"`
	Humidity                  string `json:"relative_humidity_2m"`
	Precipitation             string `json:"precipitation"`
	Precipitation_Probability string `json:"precipitation_probability"`
	WeatherCode               string `json:"weather_code"`
	WindDirection             string `json:"wind_direction_10m"`
	ApparentTemperature       string `json:"apparent_temperature"`
	CloudCover                string `json:"cloud_cover"`
}

type HourlyUnits struct {
	Time                      string `json:"time"`
	Interval                  string `json:"interval"`
	Temp                      string `json:"temperature_2m"`
	WindSpeed                 string `json:"wind_speed_10m"`
	Humidity                  string `json:"relative_humidity_2m"`
	Precipitation             string `json:"precipitation"`
	Precipitation_Probability string `json:"precipitation_probability"`
	WeatherCode               string `json:"weather_code"`
}

type Hourly struct {
	Time                      []string  `json:"time"`
	Temp2M                    []float64 `json:"temperature_2m"`
	WindSpeed                 []float64 `json:"wind_speed_10m"`
	Humidity                  []int64   `json:"relative_humidity_2m"`
	Precipitation             []float64 `json:"precipitation"`
	Precipitation_Probability []int     `json:"precipitation_probability"`
	WeatherCode               []int     `json:"weather_code"`
}

type Item struct {
	CurrentWeather CurrentWeather `json:"current"`
	CurrentUnits   CurrentUnits   `json:"current_units"`
	Hourly         Hourly         `json:"hourly"`
	HourlyUnits    HourlyUnits    `json:"hourly_units"`
}

var cache sync.Map
var lastFetch time.Time
var cacheExpiration = 5 * time.Minute

func makeDataReq() (error, Item) {

	//make get request for data with parameters

	baseURL := "https://api.open-meteo.com/v1/forecast"

	params := url.Values{}
	params.Add("latitude", "-29.3167")
	params.Add("longitude", "27.4833")
	params.Add("current", "temperature_2m,wind_speed_10m,relative_humidity_2m,precipitation,precipitation_probability,weather_code,wind_direction_10m,apparent_temperature,cloud_cover")
	params.Add("hourly", "temperature_2m,wind_speed_10m,relative_humidity_2m,precipitation,precipitation_probability,weather_code")

	queryString := params.Encode()

	url := baseURL + "?" + queryString

	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		fmt.Print(err)
		return err, Item{}
	}

	client := &http.Client{}

	//send request
	resp, err := client.Do(req)

	if err != nil {
		fmt.Print(err)
		return err, Item{}
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		fmt.Print(err)
		return err, Item{}
	}

	var data Item

	error := json.Unmarshal(body, &data)

	if error != nil {
		return error, Item{}
	}

	return nil, data

}

func endpoints() {
	router := gin.Default()

	router.GET("/all", getWeather)

	router.GET("/hourly", getHourly)

	router.GET("/current", currentWeather)

	go func() {
		for {
			time.Sleep(cacheExpiration)
			// Trigger background refresh of cache
			err, _ := makeDataReq()
			if err != nil {
				fmt.Println("Error refreshing cache:", err)
			}
		}
	}()

	router.Run("0.0.0.0:9000")

}

func getWeather(c *gin.Context) {

	now := time.Now()

	//checks for data in the cache and if time rules are respected
	if data, ok := cache.Load("data"); ok && now.Sub(lastFetch) < cacheExpiration {
		c.IndentedJSON(http.StatusOK, data)
		fmt.Println("Data fetched from cache")
		return
	}

	err, data := makeDataReq()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	cache.Store("data", data)
	lastFetch = now
	c.JSON(http.StatusOK, data)
	fmt.Println("Data fetched from remote API and stored in cache")
}

func getHourly(c *gin.Context) {
	mycache := &cache

	data, ok := mycache.Load("data")
	if !ok {
		return
	}

	var x = data
	value, ok := x.(Item)

	if !ok {
		return
	}

	hourlyData := value.Hourly

	c.IndentedJSON(http.StatusOK, hourlyData)

}

func currentWeather(c *gin.Context) {
	mycache := &cache

	data, ok := mycache.Load("data")
	if !ok {
		return
	}

	var x = data
	value, ok := x.(Item)

	if !ok {
		return
	}

	currentWthData := value.CurrentWeather
	units := value.CurrentUnits

	c.IndentedJSON(http.StatusOK, currentWthData)
	c.IndentedJSON(http.StatusOK, units)

}
