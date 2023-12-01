// TODO
// Hide API key input
// Cleaner error handling

package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	appKey, defaultLocation := getConfiguration() // Call function to ensure config file exists

	if len(os.Args) > 1 { // check for user inut
		switch os.Args[1] {
		case "--help": // Show available commands
			normal := "weather - gets the forecast for the default location (will be Bristol if not configured)"
			setDef := "set-default {arg} - Allows you to set a new loction as the default"
			setKey := "set-key {arg} - Allows you to add/replace a new applictaion key"
			diffLocation := "weather {arg} - Gets the forecast for a location whatever city/town is specified as the arg"

			fmt.Printf(
				"Available arguments:\n%s\n%s\n%s\n%s",
				normal, setDef, setKey, diffLocation,
			)
			os.Exit(0)

		case "set-default": // take input to replace default location in json file
			newDefault := input("Enter the name of a UK city to set as your default: ")
			err := updateDefaultCity(newDefault)
			handleError(err)
			os.Exit(0)

		case "set-key": // take input to replace default input in json file
			newAppKey := input("Please enter your Met Office Application Key (visible): ")
			err := updateAppKey(newAppKey)
			handleError(err)
			os.Exit(0)

		default: // Set
			defaultLocation = os.Args[1]
		}
	}

	if appKey == "" { // Ask for app key when field is empty in config file
		appKey = input("Please enter your Met Office Application Key: ")
		configData := &Config{
			ApplicationKey: appKey,
			DefaultCity:    defaultLocation,
		}
		err := writeConfig(configData) // Update key held in config file
		handleError(err)
	}

	areaId := getSiteId(defaultLocation, appKey)         // Return location ID from Met API
	date, futureForecasts := getForecast(areaId, appKey) // Return a map of today's forasts form the Met API

	colorCity := formatColor(defaultLocation, "") // Format terminal output
	separatorLine := strings.Repeat("_", 30)
	borderIcon := weatherIcons["Border"]
	headerMsg := fmt.Sprintf(
		"\n\n%s Met Office forecast for %s: %s %s\n",
		borderIcon, colorCity, date, borderIcon,
	)
	fmt.Println(headerMsg)

	for _, forecast := range futureForecasts { // Format elements from the returned forecast map

		forcastTime, err := strconv.Atoi(forecast["$"].(string))
		handleError(err)

		timeHours := forcastTime / 60 // time needs to be counted in minutes to match the format returned by the Met API
		timeString := strconv.Itoa(timeHours) + ":00:00"
		timeIcon := timeIcons[timeString]
		timeFinal := fmt.Sprintf("| Time: %s %s", timeString, timeIcon)

		weatherCode := forecast["W"].(string)
		weatherInfo, ok := metWeatherCodes[weatherCode] // Match code returned from API to it's weather value
		if !ok {
			fmt.Printf("Invalid weather code: %s", forecast["W"].(string))
			return
		}

		// Format the other relevant parts of the map and print to the console
		weatherType := fmt.Sprintf("| Weather type: %s %s", weatherInfo.Description, weatherInfo.Icon)

		tempVal := getColorEncoded(forecast["T"].(string), 15, "C")
		temperature := fmt.Sprintf("| Temperature: %s", tempVal)

		feelsLikeVal := getColorEncoded(forecast["F"].(string), 15, "C")
		feelsLike := fmt.Sprintf("| Feels like: %s", feelsLikeVal)

		windSpeedVal := getColorEncoded(forecast["S"].(string), 15, "mph")
		windSpeed := fmt.Sprintf("| Wind speed: %s", windSpeedVal)

		rainChanceVal := getColorEncoded(forecast["Pp"].(string), 49, "%")
		rainChance := fmt.Sprintf("| Chance of rain: %s", rainChanceVal)

		message := fmt.Sprintf(
			"%s\n\n%s\n%s\n%s\n%s\n%s\n%s\n\n",
			separatorLine, timeFinal, weatherType, temperature,
			feelsLike, windSpeed, rainChance,
		)
		fmt.Printf(message)
	}

	println(separatorLine)

}
