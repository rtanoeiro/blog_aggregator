package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	CurrentUserName string `json:"current_user"`
	DBURL           string `json:"db_url"`
}

func GetHomeDir() string {
	homeDir, _ := os.UserHomeDir()
	return homeDir
}

func ReadConfigFile(configPath string) (Config, error) {
	fileBytes, fileError := os.ReadFile(GetHomeDir() + configPath)

	if fileError != nil {
		fmt.Println("Error reading file from", configPath, "Error:", fileError)
		return Config{}, fileError
	}

	config := Config{}

	errUM := json.Unmarshal(fileBytes, &config)

	if errUM != nil {
		fmt.Println("Got error Unmarshling Data. Error:", errUM)
		return Config{}, errUM
	}

	return config, fileError
}

func (config *Config) SetUser(username string) error {

	config.CurrentUserName = username

	configJSON, errorMarshal := json.Marshal(config)

	if errorMarshal != nil {
		fmt.Println("Unable to marshal JSON data. Error:", errorMarshal)
	}

	os.WriteFile(GetHomeDir()+configFile, configJSON, 0644)
	return nil
}

func (config *Config) GetCurrentUser() string {
	return config.CurrentUserName
}

func (config *Config) GetUserConfig(username string) error {

	dbURL := config.DBURL
	fmt.Println("User Postgres URL:", dbURL)
	return nil
}
