package main

import (
	"encoding/json"
	"fmt"
	"strconv"

	"net/http"
	"os"

	"github.com/AlecAivazis/survey/v2"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	mapBoxAPIkey := os.Getenv("MAPBOX_API_KEY")
	openWeatherMapAPIkey := os.Getenv("OPENWEATHERMAP_API_KEY")
	var city string
	fmt.Println("What city do you want the weather for?")
	fmt.Scan(&city)
	var mapBoxEndpoint = "https://api.mapbox.com/geocoding/v5/mapbox.places/" + city + ".json?access_token=" + mapBoxAPIkey
	resp, err := http.Get(mapBoxEndpoint)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	type Message struct {
		Features []struct {
			PlaceName string    `json:"place_name"`
			Center    []float64 `json:"center"`
		}
	}
	var m Message
	json.NewDecoder(resp.Body).Decode(&m)

	options := []string{}

	for _, c := range m.Features {
		options = append(options, c.PlaceName)
	}

	selection := ""
	prompt := &survey.Select{
		Message: "Choose the correct location:",
		Options: options,
	}
	survey.AskOne(prompt, &selection)
	index := 0
	for i, c := range m.Features {
		if c.PlaceName == selection {
			index = i
			break
		}
	}

	lat := strconv.FormatFloat(m.Features[index].Center[0], 'f', -1, 64)
	long := strconv.FormatFloat(m.Features[index].Center[1], 'f', -1, 64)

	openweathermap := "https://api.openweathermap.org/data/2.5/onecall?lat=" + lat + "&lon=" + long + "&exclude=minutely,hourly,alerts&appid=" + openWeatherMapAPIkey

	resp2, err := http.Get(openweathermap)
	if err != nil {
		panic(err)
	}
	defer resp2.Body.Close()

	type Weather struct {
		Current struct {
			Temp    float64 `json:"temp"`
			Weather []struct {
				Description string `json:"description"`
			}
		}
	}

	var n Weather
	json.NewDecoder(resp2.Body).Decode(&n)
	fmt.Println(n.Current.Temp)
	fmt.Println(n.Current.Weather[0])
}
