package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"test/model"
)

// GetConfig build application config form file on disk
func GetConfig() *model.Config {
	file, fileError := os.Open("config.json")
	defer file.Close()

	if fileError != nil {
		fmt.Println("Error occurred while reading config file", fileError)
		panic(fileError)
	}

	decoder := json.NewDecoder(file)
	config := model.Config{}
	decodeError := decoder.Decode(&config)

	if decodeError != nil {
		fmt.Println("Error occurred when decoding boundary file", decodeError)
		panic(fileError)
	}

	return &config
}
