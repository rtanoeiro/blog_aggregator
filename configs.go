package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config map[string]struct {
	DBURL string `json:"db_url"`
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

	config := make(Config)

	errUM := json.Unmarshal(fileBytes, &config)

	if errUM != nil {
		fmt.Println("Got error Unmarshling Data. Error:", errUM)
		return Config{}, errUM
	}

	return config, fileError
}

func (config *Config) AddUser(username string) error {

	_, ok := (*config)[username]

	if ok {
		return fmt.Errorf("%s already present in config, skipping it", username)
	}

	temp := (*config)[username]
	temp.DBURL = "postgres://example"
	(*config)[username] = temp

	configJSON, errorMarshal := json.Marshal(config)

	if errorMarshal != nil {
		fmt.Println("Unable to marshal JSON data. Error:", errorMarshal)
	}

	os.WriteFile(GetHomeDir()+configFile, configJSON, 0644)
	return nil
}

func (config *Config) GetUserConfig(username string) error {

	dbURL, ok := (*config)[username]

	if !ok {
		return fmt.Errorf("user %s not present in config file, unable to retrieve data", username)
	}

	fmt.Println(username, "already present in config.")
	fmt.Println("Postgres URL:", dbURL)
	return nil
}
