package cache

import (
	"fmt"
	"main/logger"
	"os"

	"gopkg.in/yaml.v2"
)

type Data struct {
	Ids      map[int][]int `yaml:"ids"`
	Tappings []map[int]int `yaml:"tappings"`
}

func ReadYAMLFile(filename string) *Data {
	isFileExist(filename)

	data, err := os.ReadFile(filename)
	if err != nil {
		logger.Logger.Println(err.Error())
		return nil
	}

	var config Data
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		logger.Logger.Println(err.Error())
		os.WriteFile(filename, nil, 0644)
		return &config
	}

	return &config
}

func WriteYAMLFile(filename string, ids map[int][]int, tappings []map[int]int) {
	isFileExist(filename)

	var config Data
	data, err := os.ReadFile(filename)
	if err != nil {
		logger.Logger.Println(err.Error())
		return
	}

	err = yaml.Unmarshal(data, &config)
	if err != nil {
		logger.Logger.Println(err.Error())
		return
	}

	if len(ids) > 0 {
		config.Ids = nil
		config.Ids = ids
	}

	if len(tappings) > 0 {
		config.Tappings = nil
		config.Tappings = append(config.Tappings, tappings...)
	} else {
		config.Tappings = nil
	}

	yamlData, err := yaml.Marshal(&config)
	if err != nil {
		logger.Logger.Println(err.Error())
		return
	}

	err = os.WriteFile(filename, yamlData, 0644)
	if err != nil {
		logger.Logger.Println(err.Error())
		return
	}
}

func isFileExist(filename string) error {
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		config := &Data{}

		yamlData, err := yaml.Marshal(config)
		if err != nil {
			logger.Logger.Println(err.Error())
			return err
		}

		err = os.WriteFile(filename, yamlData, 0644)
		if err != nil {
			logger.Logger.Println(err.Error())
			return err
		}
	}

	return nil
}

func IdExists(config *Data, id int) bool {
	for _, existingIDs := range config.Ids {
		for _, existingID := range existingIDs {
			if existingID == id {
				return true
			}
		}
	}
	return false
}

func DeleteIds(filename string, yaml *Data) *Data {
	ids := make(map[int][]int)
	if yaml != nil {
		yaml.Ids = ids
		WriteYAMLFile(filename, ids, yaml.Tappings)
	}
	return yaml
}

func (y *Data) ReplaceId(oldId, newId int) error {
	for i, existingIDs := range y.Ids {
		for j, id := range existingIDs {
			if id == oldId {
				y.Ids[i][j] = newId
			}
		}
	}
	return nil
}

func TappingIdExists(config *Data, id int) bool {
	for _, values := range config.Tappings {
		for key := range values {
			if key == id {
				return true
			}
		}
	}
	return false
}

func FindTappingIdValue(config *Data, id int) int {
	for _, values := range config.Tappings {
		for key, v := range values {
			if key == id {
				return v
			}
		}
	}
	return 0
}

func UpdateTappingValue(config *Data, id, newValue int) error {
	for _, values := range config.Tappings {
		for key := range values {
			if key == id {
				values[id] = newValue
				return nil
			}
		}
	}
	logger.Logger.Printf("tapping with ID %d not found", id)
	return fmt.Errorf("tapping with ID %d not found", id)
}
