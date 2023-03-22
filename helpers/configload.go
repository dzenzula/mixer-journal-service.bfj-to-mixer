package helpers

import (
	"encoding/json"
	"main/models"
	"os"
)

func LoadAPIConfig(filePath string) *models.ConfigAPI {
	file, err := os.Open(filePath)
	if err != nil {
		return nil
	}
	defer file.Close()

	config := &models.ConfigAPI{}
	decoder := json.NewDecoder(file)
	err = decoder.Decode(config)
	if err != nil {
		return nil
	}

	return config
}

func LoadPathConfig(filePath string) *models.ConfigPath {
	file, err := os.Open(filePath)
	if err != nil {
		return nil
	}
	defer file.Close()

	config := &models.ConfigPath{}
	decoder := json.NewDecoder(file)
	err = decoder.Decode(config)
	if err != nil {
		return nil
	}

	return config
}
