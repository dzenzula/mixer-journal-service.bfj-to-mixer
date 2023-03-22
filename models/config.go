package models

type ConfigAPI struct {
	ApiUrl             string `json:"apiUrl"`
	ApiGetLastJournals string `json:"apiGetLastJournals"`
	ApiGetjournal      string `json:"apiGetJournal"`
	ApiPostAuth        string `json:"apiPostAuth"`
	ApiGetListBF       string `json:"apiGetListBF"`
}

type ConfigPath struct {
	ConfigAPIPath string `json:"configAPIPath"`
	CachePath     string `json:"cachePath"`
	AuthPath      string `json:"authPath"`
}
