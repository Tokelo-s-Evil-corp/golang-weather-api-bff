# Weather API in Go with Sync Map
This document outlines a Go server exposing a REST endpoint for retrieving current weather data. It utilizes a sync map for efficient data storage and retrieval.

## User Interface (Android App)
[Here](https://github.com/Tokelo-s-Evil-corp/weather-ui-mobile)

## Dependencies:
- GIN
- encoding/json
- fmt
- io
- net/http
- net/url
- sync
- time
  
Data Source:

The server fetches weather data from [open-meteo.com](https://opeb-meteo.com).


## Gin server log

![](https://github.com/Tokelo-s-Evil-corp/weather-ui-mobile/blob/main/Weather-golang-server.png)


Technology Stack:

## Go programming language
Sync Map (concurrent map) for fast data access
Implementation
## 1. Weather struct:

Go
```go
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
```

## 2. Sync Map:

Go
```go
var cache sync.Map
```

## 3. Fetch weather data:

This function retrieves weather data from the API and stores it in the sync map.

Go
```go
func getWeather(c *gin.Context) {

	now := time.Now()

	

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
```


## 4. Gin Endpoints 

```go
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

	router.Run("localhost:8001")

}
```
Use code with caution.

## 5. Main function:

Go

```go
package main

func main() {

	endpoints()
}
```


## Running the Server:


- Install dependencies (go get "github.com/gin-gonic/gin").
Run the server (go run .).
Accessing the API:

Once the server is running, you can access the current weather data using:

curl http://localhost:8001/current
This will return the weather data in JSON format.

Benefits of Sync Map:

Thread-safety: Provides safe access to data from concurrent routines.
Fast lookups: Enables efficient data retrieval with O(1) average-case complexity.
This improved performance allows the server to handle frequent requests efficiently.
