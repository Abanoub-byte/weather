package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
)

type apiConfigData struct {
	OpenWeatherMapApiKey string `json:"OpenWeatherMapApiKey"`
}

type weatherData struct {
	Name string `json:"name"`
	Main struct {
		Kelvin float64 `json:"temp"`
	} `json:"main"`
}

func loadApiConfig(filename string) (apiConfigData, error) {
	bytes, err := os.ReadFile(filename)

	if err != nil {
		return apiConfigData{}, err
	}

	var c apiConfigData

	err = json.Unmarshal(bytes, &c)

	if err != nil {
		return apiConfigData{}, err
	}

	return c, nil
}

func hello(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		http.Error(res, "method not allowed ", http.StatusMethodNotAllowed)
		return
	}

	_, err := res.Write([]byte("hello from go"))
	if err != nil {
		http.Error(res, "internal server error", http.StatusInternalServerError)
		return
	}
	res.Write([]byte("hello from go"))
}

func query(city string) (weatherData, error) {

	apiConfig, err := loadApiConfig(".apiConfig")
	if err != nil {
		return weatherData{}, err
	}

	resp, err := http.Get("http://api.openweathermap.org/data/2.5/weather?APPID=" + apiConfig.OpenWeatherMapApiKey + "&q=" + city)
	if err != nil {
		return weatherData{}, err
	}

	defer resp.Body.Close()
	var d weatherData
	if err := json.NewDecoder(resp.Body).Decode(&d); err != nil {
		return weatherData{}, err
	}

	return d, nil
}

func getWeather(res http.ResponseWriter, req *http.Request) {
	city := strings.SplitN(req.URL.Path, "/", 3)[2]
	data, err := query(city)

	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	res.Header().Set("Content-Type", "application/json; charset=utf-8")

	json.NewEncoder(res).Encode(data)

}

func main() {

	http.HandleFunc("/hello", hello)
	http.HandleFunc("/weather/", getWeather)

	err := http.ListenAndServe(":8081", nil)
	if err != nil {
		fmt.Println("error starting the server ", err)
	}
}
