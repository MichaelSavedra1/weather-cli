// TODO
// Split main.go file into smaller scripts
// Get icons installed
// Hide API key input
// Cleaner error handling

package main

import (
	"encoding/json"
	"fmt"
	"github.com/fatih/color"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"
)

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
	configFileName     = "config.json"
	maxRetries         = 10
	retryInterval      = 2 * time.Second
)

type Config struct {
	ApplicationKey string `json:"application_key"`
	DefaultCity    string `json:"default_city"`
}

func updateDefaultCity(defaultCity string) error {
	configData, err := readConfig()
	if err != nil {
		return err
	}

	configData.DefaultCity = defaultCity

	return writeConfig(configData)
}

func updateAppKey(appKey string) error {
	configData, err := readConfig()
	if err != nil {
		return err
	}

	configData.ApplicationKey = appKey

	return writeConfig(configData)
}

func readConfig() (*Config, error) {
	if _, err := os.Stat(configFileName); os.IsNotExist(err) {
		// File does not exist, create it with default values
		configData := &Config{
			ApplicationKey: "",
			DefaultCity:    "Bristol",
		}
		err := writeConfig(configData)
		if err != nil {
			return nil, err
		}
	}

	file, err := ioutil.ReadFile(configFileName)
	if err != nil {
		return nil, err
	}

	var configData Config
	err = json.Unmarshal(file, &configData)
	if err != nil {
		return nil, err
	}

	return &configData, nil
}

func writeConfig(configData *Config) error {
	configJSON, err := json.MarshalIndent(configData, "", "    ")
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(configFileName, configJSON, 0644)
	if err != nil {
		return err
	}

	return nil
}

func getConfiguration() (string, string) {
	configData, err := readConfig()
	handleError(err)

	return configData.ApplicationKey, configData.DefaultCity
}

func input(prompt string) string {
	fmt.Print(prompt)
	var input string
	fmt.Scanln(&input)
	return input
}

func handleError(err error) {
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}

func getRequest(url string, maxRetries int, retryInterval time.Duration) (int, []byte, error) {
	// Retry strategy enabled
	for i := 0; i < maxRetries; i++ {

		// Make HTTP request
		res, err := http.Get(url)
		if err != nil {
			continue // Retry on error
		}

		// Read response body
		body, err := ioutil.ReadAll(res.Body)
		res.Body.Close()

		// Check HTTP status code
		if res.StatusCode >= 500 {
			fmt.Printf("Received status code %d, retrying...\n", res.StatusCode)
			time.Sleep(retryInterval)
			continue // Retry on 5xx status codes
		}

		// Process successful
		return res.StatusCode, body, nil
	}

	return 0, nil, fmt.Errorf("Max retries reached")
}

func getSiteId(city string, appKey string) string {
	siteListURL := fmt.Sprintf("%s/val/wxfcs/all/%s/sitelist?key=%s", baseURL, requiredDataType, appKey)
	resStatus, resBody, err := getRequest(siteListURL, maxRetries, retryInterval)
	handleError(err)

	if resStatus != 200 {
		fmt.Println("There was an error connecting to the Met Office API.")
		os.Exit(1)
	}

	// Get the text object from web response
	responseText := string(resBody)

	// Initialise a map to use for converting the text into json
	var responseData map[string]interface{}

	// Get json from text onject
	err = json.Unmarshal([]byte(responseText), &responseData)
	handleError(err)

	// responseData should now hold the JSON
	locationsList := responseData["Locations"].(map[string]interface{})["Location"].([]interface{})

	areaId := ""

	for _, location := range locationsList {
		loc := location.(map[string]interface{})
		if loc["name"].(string) == city {
			areaId = loc["id"].(string)
			break
		}
	}
	if areaId == "" {
		fmt.Printf("No area ID found for city %s.", city)
		os.Exit(1)
	}
	return areaId
}

func main() {
	appKey, defaultCity := getConfiguration()

	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "--help":
			fmt.Println("Help options.....")
			os.Exit(0)

		case "set-default":
			newDefault := input("Enter the name of a UK city to set as your default: ")
			err := updateDefaultCity(newDefault)
			handleError(err)
			os.Exit(0)

		case "set-key":
			newAppKey := input("Please enter your Met Office Application Key: ")
			err := updateAppKey(newAppKey)
			handleError(err)
			os.Exit(0)

		default:
			defaultCity = os.Args[1]
		}
	}

	if appKey == "" {
		appKey = input("Please enter your Met Office Application Key: ")
		configData := &Config{
			ApplicationKey: appKey,
			DefaultCity:    defaultCity,
		}
		err := writeConfig(configData)
		handleError(err)
	}

	// Script execution to show forecast starts here....
	areaId := getSiteId(defaultCity, appKey)

	forcastEndpoint := fmt.Sprintf("%s/val/wxfcs/all/%s/%s?res=3hourly&key=%s", baseURL, requiredDataType, areaId, appKey)
	resStatus, resBody, err := getRequest(forcastEndpoint, maxRetries, retryInterval)
	handleError(err)

	if resStatus != 200 {
		fmt.Printf("An error ocurred when getting today's forcast: %s", resStatus)
		os.Exit(1)
	}

	responseText := string(resBody)

	var responseData map[string]interface{}

	err = json.Unmarshal([]byte(responseText), &responseData)
	if err != nil {
		fmt.Println("The Met office API failed to return a valid response. Please try again later.")
		os.Exit(1)
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
		handleError(err)

		if intFtime > timeNowMinutes {
			futureForecasts = append(futureForecasts, f)
		}
	}

	if len(futureForecasts) == 0 {
		// If no future forecasts, get the most recent entry (9 PM)
		futureForecasts = append(futureForecasts, forcastDataToday["Rep"].([]interface{})[len(forcastDataToday["Rep"].([]interface{}))-1].(map[string]interface{}))
	}

	fmt.Printf("\n\nMet Office Weather forecast for %s - %s\n\n", defaultCity, date)
	for _, forecast := range futureForecasts {

		forcastTime, err := strconv.Atoi(forecast["$"].(string))
		handleError(err)

		timeHours := forcastTime / 60
		timeFinal := fmt.Sprintf("Time: %s:00:00", strconv.Itoa(timeHours))

		weatherType := fmt.Sprintf("Weather type: %s", metWeatherCodes[forecast["W"].(string)])

		temperature := fmt.Sprintf("Temperature: %sC", forecast["T"].(string))

		feelsLike := fmt.Sprintf("Feels like: %sC", forecast["F"].(string))

		windSpeed := fmt.Sprintf("Wind speed: %smph", forecast["S"].(string))

		rainChance := fmt.Sprintf("Chance of rain: %s%%", forecast["Pp"].(string))

		message := fmt.Sprintf("%s\n%s\n%s\n%s\n%s\n%s\n\n", timeFinal, weatherType, temperature, feelsLike, windSpeed, rainChance)

		rainChanceInt, err := strconv.Atoi(forecast["Pp"].(string))
		handleError(err)

		if rainChanceInt > 50 {
			color.Red(message)
		} else if rainChanceInt > 25 {
			color.Yellow(message)
		} else {
			color.Cyan(message)
		}
	}

}
