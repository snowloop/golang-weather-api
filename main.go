package main

import (
	"fmt"
		"os"
		"encoding/json"
		"net/http"
		"strings"		
)

type apiConfigData struct {
	OpenWeatherApiKey string `json:"OpenWeatherApiKey`
}

type weatherData struct {
	Coord struct {
		Lon float64 `json:"lon`
		Lat float64 `json:"lat`
	}`json:"coord`
	Wind struct {
		Speed float64 `json:"speed`
		Deg float64 `json:"deg`
	}`json:"wind`
	Name string `json:"name`
	Main struct {
		Temp float64 `json:"temp`
		Feels_Like float64 `json:"feels_like`
        Temp_Min float64 `json:"temp_min`
        Temp_Max float64 `json:"temp_max`
        Pressure float64`json:"pressure`
        Humidity float64 `json:"humidity`
	}`json:"main`

}

func loadApiConfig(filename string)(apiConfigData, error){
	bytes, err := os.ReadFile(filename)
	
	if err != nil{
		return apiConfigData{}, err
	}
	var c apiConfigData
	err = json.Unmarshal(bytes,&c)

	if err != nil {
		return apiConfigData{}, err
	}

	return c,nil
}


func main(){
	http.HandleFunc("/hello",hello)
	http.HandleFunc("/weather/",func(w http.ResponseWriter, r *http.Request){
		fmt.Println(r.URL.Path)
		city := strings.SplitN(r.URL.Path,"/",3)[2]
		fmt.Println(city)
		data, err := query(city)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		json.NewEncoder(w).Encode(data)
	})

	http.ListenAndServe(":8080",nil)
}

func hello(w http.ResponseWriter, r *http.Request){
	w.Write([]byte("Hello from my computer"))
	w.Header().Set("Content-Type", "text/html")
}


func query(cityname string)(weatherData, error){
	apiConfigData, err := loadApiConfig(".env")
	fmt.Println(apiConfigData.OpenWeatherApiKey)
	if err != nil {
		return weatherData{}, err
	}
	resp, err := http.Get("https://api.openweathermap.org/data/2.5/weather?appid="+ apiConfigData.OpenWeatherApiKey + "&units=metric" + "&q="+ cityname)
	if err != nil {
		return weatherData{}, err
	}

	defer resp.Body.Close()

	var response weatherData
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return weatherData{},err
	}
	return response, nil
}