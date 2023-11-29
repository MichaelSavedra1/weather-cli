package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)
import "github.com/fatih/color"

var metWeatherCodes = map[string]string{
	"-1": "Trace rain",
	"0":  "Clear night",
	"1":  "Sunny day",
	"2":  "Partly cloudy (night)",
	"3":  "Partly cloudy (day)",
	"4":  "Not used",
	"5":  "Mist",
	"6":  "Fog",
	"7":  "Cloudy",
	"8":  "Overcast",
	"9":  "Light rain shower (night)",
	"10": "Light rain shower (day)",
	"11": "Drizzl",
	"12": "Light rain",
	"13": "Heavy rain shower (night)",
	"14": "Heavy rain shower (day)",
	"15": "Heavy rain",
	"16": "Sleet shower (night)",
	"17": "Sleet shower (day)",
	"18": "Sleet",
	"19": "Hail shower (night)",
	"20": "Hail shower (day)",
	"21": "Hail",
	"22": "Light snow shower (night)",
	"23": "Light snow shower (day)",
	"24": "Light snow",
	"25": "Heavy snow shower (night)",
	"26": "Heavy snow shower (day)",
	"27": "Heavy snow",
	"28": "Thunder shower (night)",
	"29": "Thunder shower (day)",
	"30": "Thunder",
}

const (
	baseURL            = "http://datapoint.metoffice.gov.uk/public/data"
	getRegionsEndpoint = "txt/wxfcs/regionalforecast/datatype/sitelist"
	requiredDataType   = "json"
	applicationKeyFile = "met_application_key.txt"
)

func readAPIKey() string {
	apiKey, err := ioutil.ReadFile("met_application_key.txt")
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(apiKey))
}

func writeAPIKey(apiKey string) error {
	return ioutil.WriteFile("met_application_key.txt", []byte(apiKey), 0644)
}

func getSiteId(city string, apiKey string) string {
	siteListURL := fmt.Sprintf("%s/val/wxfcs/all/%s/sitelist?key=%s", baseURL, requiredDataType, apiKey)
	res, err := http.Get(siteListURL)
	if err != nil {
		panic(err)
	}

	defer res.Body.Close() //close body of response

	if res.StatusCode != 200 {
		panic("There was an error connecting to the Met Office API. Please ensure your application key is correct.")
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	// Get the text object from web response
	responseText := string(body)

	// Initialise a map to use for converting the text into json
	var responseData map[string]interface{}

	// Get json from text onject
	err = json.Unmarshal([]byte(responseText), &responseData)
	if err != nil {
		panic("Failed to format JSON.")
	}

	// responseData should now hold the JSON
	locationsList := responseData["Locations"].(map[string]interface{})["Location"].([]interface{})

	var areaId string

	for _, location := range locationsList {
		loc := location.(map[string]interface{})
		if loc["name"].(string) == "Bristol" {
			areaId = loc["id"].(string)
			break
		}
	}

	if areaId == "" {
		panic("No area ID found for Bristol.")
	}

	return areaId
}

func main() {

	apiKey := readAPIKey()
	if apiKey == "" {
		fmt.Print("Please enter your Met Application key: ")
		fmt.Scanln(&apiKey)
		writeAPIKey(apiKey)
	}

	city := "Bristol" //Set Default city for if no CLI arg is provided.
	if len(os.Args) >= 2 {
		city = os.Args[1]
	}

	areaId := getSiteId(city, apiKey)

	forcastEndpoint := fmt.Sprintf("%s/val/wxfcs/all/%s/%s?res=3hourly&key=%s", baseURL, requiredDataType, areaId, apiKey)
	res, err := http.Get(forcastEndpoint)
	if err != nil {
		panic(err)
	}

	defer res.Body.Close() //close body of response

	if res.StatusCode != 200 {
		panic("An error ocurred when getting today's forcast.")
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		panic("Failed to parse JSON data for today's weather.")
	}

	responseText := string(body)
	var responseData map[string]interface{}

	err = json.Unmarshal([]byte(responseText), &responseData)
	if err != nil {
		panic("The was an error with the Met Office API.")
	}

	forcastDataToday := responseData["SiteRep"].(map[string]interface{})["DV"].(map[string]interface{})["Location"].(map[string]interface{})["Period"].([]interface{})[0].(map[string]interface{})
	date := forcastDataToday["value"].(string)

	currentUTC := time.Now().UTC()
	timeNowMinutes := currentUTC.Hour()*60 + currentUTC.Minute()

	futureForecasts := make([]map[string]interface{}, 0)

	for _, forecast := range forcastDataToday["Rep"].([]interface{}) {
		f := forecast.(map[string]interface{})

		fTime := f["$"].(string)
		intFtime, err := strconv.Atoi(fTime)
		if err != nil {
			panic(err)
		}

		if intFtime > timeNowMinutes {
			futureForecasts = append(futureForecasts, f)
		}
	}

	if len(futureForecasts) == 0 {
		// If no future forecasts, get the most recent entry (9 PM)
		futureForecasts = append(futureForecasts, forcastDataToday["Rep"].([]interface{})[len(forcastDataToday["Rep"].([]interface{}))-1].(map[string]interface{}))
	}

	fmt.Printf("Met Office Weather forecast for %s - %s\n\n", city, date)
	for _, forecast := range futureForecasts {

		forcastTime, err := strconv.Atoi(forecast["$"].(string))
		if err != nil {
			panic(err)
		}

		timeHours := forcastTime / 60
		timeFinal := fmt.Sprintf("Time: %s:00:00", strconv.Itoa(timeHours))

		weatherType := fmt.Sprintf("Weather type: %s", metWeatherCodes[forecast["W"].(string)])

		temperature := fmt.Sprintf("Temperature: %sC", forecast["T"].(string))

		feelsLike := fmt.Sprintf("Feels like: %sC", forecast["F"].(string))

		windSpeed := fmt.Sprintf("Wind speed: %smph", forecast["S"].(string))

		rainChance := fmt.Sprintf("Chance of rain: %s%%", forecast["Pp"].(string))

		message := fmt.Sprintf("%s\n%s\n%s\n%s\n%s\n%s\n\n", timeFinal, weatherType, temperature, feelsLike, windSpeed, rainChance)

		rainChanceInt, err := strconv.Atoi(forecast["Pp"].(string))
		if err != nil {
			panic(err)
		}

		if rainChanceInt > 50 {
			color.Red(message)
		} else if rainChanceInt > 25 {
			color.Yellow(message)
		} else {
			color.Cyan(message)
		}
	}

}
