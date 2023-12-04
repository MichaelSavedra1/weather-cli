package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

func getForecast(areaId string, appKey string, extended bool) ([]map[string]interface{}, error) {
	forcastEndpoint := fmt.Sprintf(
		"%s/val/wxfcs/all/%s/%s?res=3hourly&key=%s",
		baseURL, requiredDataType, areaId, appKey,
	)

	resStatus, resBody, err := getRequest(forcastEndpoint, maxRetries, retryInterval) // Get forecast data
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

	var responseArray []map[string]interface{} // this is the response onject for the function

	// Slice into the returned map from the met (full format available in the resources dir)
	allForecasts := responseData["SiteRep"].(map[string]interface{})["DV"].(map[string]interface{})["Location"].(map[string]interface{})["Period"].([]interface{})

	forecastsToday := allForecasts[0].(map[string]interface{})
	dateToday := forecastsToday["value"].(string)

	// get current time as logic to decide which forecasts form today to include
	currentUTC := time.Now().UTC()
	timeNowMinutes := currentUTC.Hour()*60 + currentUTC.Minute()

	futureForecasts := make([]map[string]interface{}, 0) // empty list for today's future forecasts

	for _, forecast := range forecastsToday["Rep"].([]interface{}) { // only append future times to current day section
		f := forecast.(map[string]interface{})
		fTime := f["$"].(string)
		intFtime, err := strconv.Atoi(fTime)
		handleError(err)

		if intFtime > timeNowMinutes {
			futureForecasts = append(futureForecasts, f)
		}
	}

	if len(futureForecasts) == 0 { // If no future forecasts, get the most recent entry (9 PM)
		futureForecasts = append(futureForecasts, forecastsToday["Rep"].([]interface{})[len(forecastsToday["Rep"].([]interface{}))-1].(map[string]interface{}))
	}

	todaysForecasts := make(map[string]interface{}) // create instance of map to add to return list
	todaysForecasts["date"] = dateToday
	todaysForecasts["forecasts"] = futureForecasts

	responseArray = append(responseArray, todaysForecasts) // add struct instance to response array object

	if extended { //add additional forecasts if requested
		for _, forecast := range allForecasts[1:] { // loop through each dict in forecast list

			additionalForecast := make(map[string]interface{}) // build dict for each additional day to store date and forecast data
			futureData := make([]map[string]interface{}, 0)    // build empty list to append to the forecast key of 'additionalForecast'
			f := forecast.(map[string]interface{})             // convert curernt item into map format
			additionalForecast["date"] = f["value"].(string)   // add date to date key

			for _, item := range f["Rep"].([]interface{}) {
				i := item.(map[string]interface{})
				futureData = append(futureData, i)
			}
			additionalForecast["forecasts"] = futureData
			// append the extended forecast to the return array
			responseArray = append(responseArray, additionalForecast)
		}

	}

	return responseArray, nil
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

	responseText := string(resBody) // Get the text object from web response

	var responseData map[string]interface{} // Initialise a map to use for converting the text into json

	err = json.Unmarshal([]byte(responseText), &responseData) // Get json from text onject
	handleError(err)

	locationsList := responseData["Locations"].(map[string]interface{})["Location"].([]interface{}) // responseData should now hold the JSON

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
