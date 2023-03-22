package cache

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Data struct {
	Ids []int `yaml:"ids"`
}

func ReadYAMLFile(fileName string) *Data {
	data, err := ioutil.ReadFile(fileName)
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

func WriteYAMLFile(filename string, config *Data, id int) {
	config.Ids = append(config.Ids, id)

	yamlData, err := yaml.Marshal(&config)
	if err != nil {
		return
	}

	err = ioutil.WriteFile(filename, yamlData, 0644)
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
