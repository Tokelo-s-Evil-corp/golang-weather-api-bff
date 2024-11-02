# Weather API in Go with Sync Map
This document outlines a Go server exposing a REST endpoint for retrieving current weather data. It utilizes a sync map for efficient data storage and retrieval.

## User Interface (Android App)
[Here](https://github.com/Tokelo-s-Evil-corp/weather-ui-mobile)

## Dependencies:

- encoding/json
- fmt
- io
- net/http
- net/url
- sync
- time
  
Data Source:

The server fetches weather data from open-mateo.com (replace with your preferred API).
Technology Stack:

## Go programming language
Sync Map (concurrent map) for fast data access
Implementation
## 1. Weather struct:

Go
```go
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

type Hourly struct {
	Times []string `json:"time"`
}

type Item struct {
	CurrentWeather CurrentWeather `json:"current"`
	CurrentUnits   CurrentUnits   `json:"current_units"`
	Hourly         Hourly         `json:"hourly"`
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

	// var data Item

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
