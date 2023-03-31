package helpers

import (
	"encoding/json"
	"os"
)

func LoadAPIConfig(filePath string) *ConfigAPI {
	file, err := os.Open(filePath)
	if err != nil {
		return nil
	}
	defer file.Close()

	config := &ConfigAPI{}
	decoder := json.NewDecoder(file)
	err = decoder.Decode(config)
	if err != nil {
		return nil
	}

	return config
}

func LoadPathConfig(filePath string) *ConfigPath {
	file, err := os.Open(filePath)
	if err != nil {
		return nil
	}
	defer file.Close()

	config := &ConfigPath{}
	decoder := json.NewDecoder(file)
	err = decoder.Decode(config)
	if err != nil {
		return nil
	}

	return config
}

func LoadAPIMixConfig(filePath string) *ConfigAPIMix {
	file, err := os.Open(filePath)
	if err != nil {
		return nil
	}
	defer file.Close()

	config := &ConfigAPIMix{}
	decoder := json.NewDecoder(file)
	err = decoder.Decode(config)
	if err != nil {
		return nil
	}

	return config
}

func LoadTimeConfig(filePath string) *ConfigTime {
	file, err := os.Open(filePath)
	if err != nil {
		return nil
	}
	defer file.Close()

	config := &ConfigTime{}
	decoder := json.NewDecoder(file)
	err = decoder.Decode(config)
	if err != nil {
		return nil
	}

	return config
}
