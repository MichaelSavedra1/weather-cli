package main

import "time"

// icons lib
var weatherIcons = map[string]string{
	"Cloud":         "\033[38;2;135;206;250m☁\033[0m",
	"Rainy Cloud":   "\033[38;2;0;0;255m🌧\033[0m",
	"Sunny Cloud":   "\033[38;2;255;255;0m⛅\033[0m",
	"Thunder Cloud": "\033[38;2;255;0;0m🌩\033[0m",
	"Snow":          "\033[38;2;255;255;255m❄\033[0m",
	"Wind":          "\033[38;2;100;149;237m🌫\033[0m",

	"Night Other": "🌚",
	"Night Clear": "🌛",
	"Border":      "🌐",
	"Heavy Rain":  "💦",
	"Light Rain":  "💧",
	"Sun":         "🌞",
	"Fog":         "💨",
	"Warning":     "🚩",
	"Ice":         "🧊",
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
	"-1": {Description: "Trace rain", Icon: weatherIcons["Light Rain"]},
	"0":  {Description: "Clear night", Icon: weatherIcons["Night Clear"]},
	"1":  {Description: "Sunny day", Icon: weatherIcons["Sun"]},
	"2":  {Description: "Partly cloudy (night)", Icon: weatherIcons["Night Other"]},
	"3":  {Description: "Partly cloudy (day)", Icon: weatherIcons["Sunny Cloud"]},

	"4": {Description: "Not used", Icon: ""},

	"5": {Description: "Mist", Icon: weatherIcons["Fog"]},
	"6": {Description: "Fog", Icon: weatherIcons["Fog"]},

	"7": {Description: "Cloudy", Icon: weatherIcons["Cloud"]},
	"8": {Description: "Overcast", Icon: weatherIcons["Cloud"]},

	"9": {Description: "Light rain shower (night)", Icon: weatherIcons["Night Other"]},

	"10": {Description: "Light rain shower (day)", Icon: weatherIcons["Light Rain"]},
	"11": {Description: "Drizzle", Icon: weatherIcons["Light Rain"]},
	"12": {Description: "Light rain", Icon: weatherIcons["Light Rain"]},

	"13": {Description: "Heavy rain shower (night)", Icon: weatherIcons["Heavy Rain"]},
	"14": {Description: "Heavy rain shower (day)", Icon: weatherIcons["Heavy Rain"]},
	"15": {Description: "Heavy rain", Icon: weatherIcons["Heavy Rain"]},

	"16": {Description: "Sleet shower (night)", Icon: weatherIcons["Heavy Rain"]},
	"17": {Description: "Sleet shower (day)", Icon: weatherIcons["Heavy Rain"]},
	"18": {Description: "Sleet", Icon: weatherIcons["Warning"]},
	"19": {Description: "Hail shower (night)", Icon: weatherIcons["Light Rain"]},
	"20": {Description: "Hail shower (day)", Icon: weatherIcons["Light Rain"]},
	"21": {Description: "Hail", Icon: weatherIcons["Warning"]},

	"22": {Description: "Light snow shower (night)", Icon: weatherIcons["Snow"]},
	"23": {Description: "Light snow shower (day)", Icon: weatherIcons["Snow"]},
	"24": {Description: "Light snow", Icon: weatherIcons["Snow"]},
	"25": {Description: "Heavy snow shower (night)", Icon: weatherIcons["Ice"]},
	"26": {Description: "Heavy snow shower (day)", Icon: weatherIcons["Snow"]},
	"27": {Description: "Heavy snow", Icon: weatherIcons["Ice"]},

	"28": {Description: "Thunder shower (night)", Icon: weatherIcons["Thunder Cloud"]},
	"29": {Description: "Thunder shower (day)", Icon: weatherIcons["Thunder Cloud"]},
	"30": {Description: "Thunder", Icon: weatherIcons["Warning"]},
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
