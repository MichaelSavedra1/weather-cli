package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

func getForecast(areaId string, appKey string) (string, []map[string]interface{}) {
	forcastEndpoint := fmt.Sprintf(
		"%s/val/wxfcs/all/%s/%s?res=3hourly&key=%s",
		baseURL, requiredDataType, areaId, appKey,
	)

	resStatus, resBody, err := getRequest(forcastEndpoint, maxRetries, retryInterval)
	handleError(err)

	if resStatus != 200 {
		fmt.Printf("An error ocurred when getting today's forcast: %d", resStatus)
		os.Exit(1)
	}

	responseText := string(resBody)

	var responseData map[string]interface{}

	err = json.Unmarshal([]byte(responseText), &responseData)
	if err != nil {
		fmt.Println("The Met office API failed to return a valid response. Please try again later.")
		os.Exit(1)
	}

	// Slice into the returned map from the met (full format available in the resources dir)
	forcastDataToday := responseData["SiteRep"].(map[string]interface{})["DV"].(map[string]interface{})["Location"].(map[string]interface{})["Period"].([]interface{})[0].(map[string]interface{})
	date := forcastDataToday["value"].(string)

	// get current time as logic to decide which forecasts form today to include
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
	return date, futureForecasts
}

func getSiteId(area string, appKey string) string {
	siteListURL := fmt.Sprintf(
		"%s/val/wxfcs/all/%s/sitelist?key=%s",
		baseURL, requiredDataType, appKey,
	)
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

	areaId := "" // Set default val for error handline

	for _, location := range locationsList {
		loc := location.(map[string]interface{})
		if strings.ToLower(loc["name"].(string)) == strings.ToLower(area) {
			areaId = loc["id"].(string)
			break
		}
	}
	if areaId == "" {
		fmt.Printf("No area ID found for city %s\n.", area)
		os.Exit(0)
	}
	return areaId
}
