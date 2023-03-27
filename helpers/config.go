package helpers

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
}

type ConfigPath struct {
	ConfigAPIPath string `json:"configAPIPath"`
	CachePath     string `json:"cachePath"`
	AuthPath      string `json:"authPath"`
}

var CfgPath *ConfigPath = LoadPathConfig("models/json/configPaths.json")
var CfgAPI *ConfigAPI = LoadAPIConfig("models/json/configAPI.json")
