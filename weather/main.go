// TODO
// Hide API key input
// Cleaner error handling

package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var (
	help       = flag.Bool("help", false, "show help message.")
	v          = flag.Bool("v", false, "show app version.")
	extended   = flag.Bool("extended", false, "extended forcast")
	setDefault = flag.Bool("set-default", false, "update default city")
	setKey     = flag.Bool("set-key", false, "update met office api key")
)

func main() {
	flag.Parse()
	appKey, defaultLocation := getConfiguration() // Call function to ensure config file exists
	var e = false                                 // not extended

	if flag.NFlag() > 1 {
		fmt.Println("`weather` takes exactly one arg - see `weather --help`")
		return
	}
	if flag.NArg() == 0 && flag.NFlag() == 1 {
		if *help {
			normal := "| 'weather' -> gets the forecast for the default location (will be Bristol if not configured)"
			setDef := "| 'weather --set-default' -> Allows you to set a new loction as the default"
			setKey := "| 'weather --set-key' -> Allows you to add/replace a new applictaion key"
			diffLocation := "| 'weather {arg}' -> Gets the forecast for a location whatever city/town is specified as the arg"
			extend := "| 'weather --extended' -> shows the next 5 days of forecast data. Can optionally use 'weather {arg} --extended'"
			fmt.Printf(
				"Available arguments:\n\n%s\n\n%s\n\n%s\n\n%s\n\n%s",
				normal, setDef, setKey, diffLocation, extend,
			)
			return
		}

		if *v {
			fmt.Println("v1.0.0")
			return
		}

		if *setDefault {
			newDefault := input("Enter the name of a UK city to set as your default: ")
			err := updateDefaultCity(newDefault)
			handleError(err)
			return
		}

		if *extended {
			e = true
			defaultLocation = os.Args[1]
		}

		if *setKey {
			newAppKey := input("Please enter your Met Office Application Key (visible): ")
			err := updateAppKey(newAppKey)
			handleError(err)
			return
		}
	}

	if len(os.Args) > 2 && !strings.Contains(os.Args[2], "extended") {
		fmt.Println("malformed request see `weather --help`")
		return

	} else if len(os.Args) > 3 {
		fmt.Println("too many args - see `weather --help`")
		return
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

	if len(os.Args) > 1 {
		defaultLocation = os.Args[1]
	}

	areaId := getSiteId(defaultLocation, appKey) // Return location ID from Met API

	forecasts, err := getForecast(areaId, appKey, e) // get forecasts
	handleError(err)

	colorCity := formatColor(defaultLocation, "") // Format terminal output
	separatorLine := strings.Repeat("_", 30)
	borderIcon := weatherIcons["Border"]

	for _, forecast := range forecasts {
		formattedDate, err := formatDate(forecast["date"].(string))
		handleError(err)

		headerMsg := fmt.Sprintf(
			"\n\n%s Met Office forecast for %s: %s %s\n",
			borderIcon, colorCity, formattedDate, borderIcon,
		)
		fmt.Println(headerMsg)

		for _, forecastData := range forecast["forecasts"].([]map[string]interface{}) {
			forcastTime, err := strconv.Atoi(forecastData["$"].(string))
			handleError(err)

			timeHours := forcastTime / 60 // time needs to be counted in minutes to match the format returned by the Met API
			timeString := strconv.Itoa(timeHours) + ":00:00"
			timeIcon := timeIcons[timeString]
			timeFinal := fmt.Sprintf("| Time: %s %s", timeString, timeIcon)

			weatherCode := forecastData["W"].(string)
			weatherInfo, ok := metWeatherCodes[weatherCode] // Match code returned from API to its weather value
			if !ok {
				fmt.Printf("Invalid weather code: %s", forecastData["W"].(string))
				return
			}

			// Format the other relevant parts of the map and print to the console
			weatherType := fmt.Sprintf("| Weather type: %s %s", weatherInfo.Description, weatherInfo.Icon)

			tempVal := getColorEncoded(forecastData["T"].(string), "C")
			temperature := fmt.Sprintf("| Temperature: %s", tempVal)

			feelsLikeVal := getColorEncoded(forecastData["F"].(string), "C")
			feelsLike := fmt.Sprintf("| Feels like: %s", feelsLikeVal)

			windSpeedVal := getColorEncoded(forecastData["S"].(string), "mph")
			windSpeed := fmt.Sprintf("| Wind speed: %s", windSpeedVal)

			rainChanceVal := getColorEncoded(forecastData["Pp"].(string), "%")
			rainChance := fmt.Sprintf("| Chance of rain: %s", rainChanceVal)

			message := fmt.Sprintf(
				"%s\n\n%s\n%s\n%s\n%s\n%s\n%s\n\n",
				separatorLine, timeFinal, weatherType, temperature,
				feelsLike, windSpeed, rainChance,
			)
			fmt.Println(message)
		}
		println(separatorLine)
	}

}
