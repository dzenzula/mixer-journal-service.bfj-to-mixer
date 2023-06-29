package models

type ChemicalDTO struct {
	//NMix       int      `json:"nMix"`
	NBlock     int      `json:"nBlock"`
	Ladle      string   `json:"ladle"`
	NumSample  int      `json:"numSample"`
	NumTaphole int      `json:"numTaphole"`
	DT         *string  `json:"dt"`
	Si         *float64 `json:"si"`
	Mn         *float64 `json:"mn"`
	S          *float64 `json:"s"`
	P          *float64 `json:"p"`
	Belong     string   `json:"belong"`
}
