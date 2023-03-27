package models

type ChemMaterial struct {
	ID             int      `json:"id"`
	IDJournal      int      `json:"idJournal"`
	NameSource     string   `json:"nameSource"`
	NumberRoute    string   `json:"numberRoute"`
	WeightMaterial int      `json:"weightMaterial"`
	Fe             float64  `json:"fe"`
	FeO            float64  `json:"feO"`
	SiO2           float64  `json:"siO2"`
	MnO            *float64 `json:"mnO,omitempty"`
	CaO            float64  `json:"caO"`
	MgO            float64  `json:"mgO"`
	Fraction       float64  `json:"fraction"`
	CaOSiO2        float64  `json:"caOsiO2"`
	AddData        string   `json:"addData"`
	DtInsert       string   `json:"dtInsert"`
	DtUpdate       string   `json:"dtUpdate,omitempty"`
	CreatedOn      string   `json:"createdOn,omitempty"`
	ModifyBy       int      `json:"modifyBy,omitempty"`
}
