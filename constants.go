package main

import "time"

// icons lib
var weatherIcons = map[string]string{
	"Cloud":         "\033[38;2;135;206;250m☁\033[0m",
	"Rainy Cloud":   "\033[38;2;0;0;255m🌧\033[0m",
	"Sun":           "🌞",
	"Sunny Cloud":   "\033[38;2;255;255;0m⛅\033[0m",
	"Thunder Cloud": "\033[38;2;255;0;0m🌩\033[0m",
	"Snow":          "\033[38;2;255;255;255m❄\033[0m",
	"Wind":          "\033[38;2;100;149;237m🌫\033[0m",
	"Border":        "🌐",
}

var timeIcons = map[string]string{
	"00:00:00": "🕛",
	"03:00:00": "🕒",
	"06:00:00": "🕕",
	"09:00:00": "🕘",
	"12:00:00": "🕛",
	"15:00:00": "🕒",
	"18:00:00": "🕕",
	"21:00:00": "🕘",
}

// map used to costruct a link between metWeatherCodes map
// and weatherIcons map
type WeatherInfo struct {
	Description string
	Icon        string
}

// Weather types as defined by the Met office
var metWeatherCodes = map[string]WeatherInfo{
	"-1": {Description: "Trace rain", Icon: weatherIcons["Rainy Cloud"]},
	"0":  {Description: "Clear night", Icon: weatherIcons["Sun"]},
	"1":  {Description: "Sunny day", Icon: weatherIcons["Sun"]},
	"2":  {Description: "Partly cloudy (night)", Icon: weatherIcons["Sunny Cloud"]},
	"3":  {Description: "Partly cloudy (day)", Icon: weatherIcons["Sunny Cloud"]},
	"4":  {Description: "Not used", Icon: ""},
	"5":  {Description: "Mist", Icon: weatherIcons["Cloud"]},
	"6":  {Description: "Fog", Icon: weatherIcons["Cloud"]},
	"7":  {Description: "Cloudy", Icon: weatherIcons["Cloud"]},
	"8":  {Description: "Overcast", Icon: weatherIcons["Cloud"]},
	"9":  {Description: "Light rain shower (night)", Icon: weatherIcons["Rainy Cloud"]},
	"10": {Description: "Light rain shower (day)", Icon: weatherIcons["Rainy Cloud"]},
	"11": {Description: "Drizzle", Icon: weatherIcons["Rainy Cloud"]},
	"12": {Description: "Light rain", Icon: weatherIcons["Rainy Cloud"]},
	"13": {Description: "Heavy rain shower (night)", Icon: weatherIcons["Rainy Cloud"]},
	"14": {Description: "Heavy rain shower (day)", Icon: weatherIcons["Rainy Cloud"]},
	"15": {Description: "Heavy rain", Icon: weatherIcons["Rainy Cloud"]},
	"16": {Description: "Sleet shower (night)", Icon: weatherIcons["Snow"]},
	"17": {Description: "Sleet shower (day)", Icon: weatherIcons["Snow"]},
	"18": {Description: "Sleet", Icon: weatherIcons["Snow"]},
	"19": {Description: "Hail shower (night)", Icon: weatherIcons["Snow"]},
	"20": {Description: "Hail shower (day)", Icon: weatherIcons["Snow"]},
	"21": {Description: "Hail", Icon: weatherIcons["Snow"]},
	"22": {Description: "Light snow shower (night)", Icon: weatherIcons["Snow"]},
	"23": {Description: "Light snow shower (day)", Icon: weatherIcons["Snow"]},
	"24": {Description: "Light snow", Icon: weatherIcons["Snow"]},
	"25": {Description: "Heavy snow shower (night)", Icon: weatherIcons["Snow"]},
	"26": {Description: "Heavy snow shower (day)", Icon: weatherIcons["Snow"]},
	"27": {Description: "Heavy snow", Icon: weatherIcons["Snow"]},
	"28": {Description: "Thunder shower (night)", Icon: weatherIcons["Thunder Cloud"]},
	"29": {Description: "Thunder shower (day)", Icon: weatherIcons["Thunder Cloud"]},
	"30": {Description: "Thunder", Icon: weatherIcons["Thunder Cloud"]},
}

// general consts
const (
	baseURL = "http://datapoint.metoffice.gov.uk/public/data"
	//getRegionsEndpoint = "txt/wxfcs/regionalforecast/datatype/sitelist"
	requiredDataType = "json"
	configFileName   = "config.json"
	maxRetries       = 10
	retryInterval    = 2 * time.Second
)
