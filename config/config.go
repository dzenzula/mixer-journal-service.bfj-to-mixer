package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type ConfigPath struct {
	CachePath string `json:"cachePath"`
}

type Config struct {
	BFJAPI ConfigAPI     `json:"bfj_api_prod"`
	MIXAPI ConfigAPIMix  `json:"mix_api_test"`
	Path   ConfigPath    `json:"paths"`
	Auth   Authorization `json:"auth"`
}

type ConfigAPI struct {
	ApiUrl               string `json:"apiUrl"`
	ApiGetLastJournals   string `json:"apiGetLastJournals"`
	ApiGetjournal        string `json:"apiGetJournal"`
	ApiPostAuthTest      string `json:"apiPostAuthTest"`
	ApiPostAuthProd      string `json:"apiPostAuthProd"`
	ApiGetListBF         string `json:"apiGetListBF"`
	ApiGetChemCoxes      string `json:"apiGetChemCoxes"`
	ApiGetChemicalsSlags string `json:"apiGetChemicalsSlags"`
	ApiGetChemMaterials  string `json:"apiGetChemMaterials"`
	ApiGetTappings       string `json:"apiGetTappings"`
}

type ConfigAPIMix struct {
	ApiUrl               string `json:"apiUrl"`
	ApiPostAuthTest      string `json:"apiPostAuthTest"`
	ApiGetLastJournals   string `json:"apiGetLastJournals"`
	ApiPostChemical      string `json:"apiPostChemical"`
	ApiPostLadleMovement string `json:"apiPostLadleMovement"`
	ApiGetLadleMovement  string `json:"apiGetLadleMovement"`
}

type Authorization struct {
	Login    string `json:"login"`
	Password string `json:"password"`
	Domain   string `json:"domain"`
}

var GlobalConfig *Config = getConfig()

func getConfig() *Config {
	var config_files [3]string
	var file *os.File
	config_name := "config.json"
	config_files[0] = "models/json/" + config_name
	config_files[1] = "config/" + config_name
	config_files[2] = config_name

	for i := 0; i < len(config_files); i++ {
		if _, err := os.Stat(config_files[i]); err == nil {
			conf_file, err := os.Open(config_files[i])
			if err != nil {
				fmt.Printf("Error opening config file: %s", err)
			}
			file = conf_file
			defer file.Close()
		}
	}

	config := &Config{}
	decoder := json.NewDecoder(file)
	err := decoder.Decode(config)
	if err != nil {
		fmt.Printf("error decoding config file: %s", err)
	}

	return config
}
