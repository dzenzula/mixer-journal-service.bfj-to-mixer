package models

type ChemicalSlag struct {
	IDJournal  int     `json:"id_journal"`
	NumTapping int     `json:"numTapping"`
	SiO2       float64 `json:"siO2"`
	Al2O3      float64 `json:"al2O3"`
	CaO        float64 `json:"caO"`
	MnO        float64 `json:"mnO"`
	MgO        float64 `json:"mgO"`
	FeO        float64 `json:"feO"`
	SS         float64 `json:"ss"`
	Och1       float64 `json:"och1"`
	Och2       float64 `json:"och2"`
}
