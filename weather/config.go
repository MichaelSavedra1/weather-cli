package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type Config struct { // Construct format for json config file
	ApplicationKey string `json:"application_key"`
	DefaultCity    string `json:"default_city"`
}

func updateDefaultCity(defaultCity string) error {
	// Update default city location held in json file
	configData, err := readConfig()
	if err != nil {
		return err
	}

	configData.DefaultCity = defaultCity

	return writeConfig(configData)
}

func updateAppKey(appKey string) error {
	// Update the app key value stored in the config file
	configData, err := readConfig()
	if err != nil {
		return err
	}

	configData.ApplicationKey = appKey

	return writeConfig(configData)
}

func readConfig() (*Config, error) {
	// 1. Create config file with default values if it is not present in path
	// 2. Return the values of the config file
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
	// Write the config to the json file
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
	// return the data from the config file
	configData, err := readConfig()
	handleError(err)

	return configData.ApplicationKey, configData.DefaultCity
}
