package cache

import (
	"os"

	"gopkg.in/yaml.v2"
)

type Data struct {
	Ids []int `yaml:"ids"`
}

func ReadYAMLFile(fileName string) *Data {
	data, err := os.ReadFile(fileName)
	if err != nil {
		return nil
	}

	var config Data
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil
	}

	return &config
}

func WriteYAMLFile(filename string, config *Data, ids map[int][]int) {
	for _, values := range ids {
		config.Ids = append(config.Ids, values...)
	}

	yamlData, err := yaml.Marshal(&config)
	if err != nil {
		return
	}

	err = os.WriteFile(filename, yamlData, 0644)
	if err != nil {
		return
	}

}

func IdExists(config *Data, id int) bool {
	for _, existingID := range config.Ids {
		if existingID == id {
			return true
		}
	}
	return false
}

func DeleteIds(filename string, yaml *Data) *Data {
	ids := make(map[int][]int)
	if yaml != nil {
		yaml.Ids = []int{}
		WriteYAMLFile(filename, yaml, ids)
	}
	return yaml
}

func (y *Data) ReplaceId(oldId, newId int) error {
	for i, id := range y.Ids {
		if id == oldId {
			y.Ids[i] = newId
		}
	}
	return nil
}
