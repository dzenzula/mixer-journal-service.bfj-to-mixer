package helpers

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	BFJAPI ConfigAPI     `json:"bfj_api_prod"`
	MIXAPI ConfigAPIMix  `json:"mix_api_test"`
	Path   ConfigPath    `json:"paths"`
	Time   ConfigTime    `json:"time"`
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

type ConfigPath struct {
	CachePath string `json:"cachePath"`
}

type ConfigAPIMix struct {
	ApiUrl          string `json:"apiUrl"`
	ApiPostAuthTest string `json:"apiPostAuthTest"`
}

type ConfigTime struct {
	OneMinuteInterval  string `json:"one_minute_interval"`
	FiveMinuteInterval string `json:"five_minute_interval"`
	TenMinuteInterval  string `json:"ten_minute_interval"`
	OneHourInterval    string `json:"one_hour_interval"`
	FiveHourInterval   string `json:"five_hour_interval"`
}

type Authorization struct {
	Login    string `json:"login"`
	Password string `json:"password"`
	Domain   string `json:"domain"`
}

var GlobalConfig *Config = getConfig()

func getConfig() *Config {
	file, err := os.Open("models/json/config.json")
	if err != nil {
		fmt.Printf("Error opening config file: %s", err)
	}
	defer file.Close()

	config := &Config{}
	decoder := json.NewDecoder(file)
	err = decoder.Decode(config)
	if err != nil {
		fmt.Printf("error decoding config file: %s", err)
	}

	return config
}
