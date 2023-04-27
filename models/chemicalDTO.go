package models

type ChemicalDTO struct {
	NMix       int     `json:"nMix"`
	Ladle      string  `json:"ladle"`
	Proba      int     `json:"proba"`
	NumTaphole int     `json:"numTaphole"`
	DT         *string `json:"dt"`
	Si         int     `json:"si"`
	Mn         int     `json:"mn"`
	S          int     `json:"s"`
	P          int     `json:"p"`
	Belong     string  `json:"belong"`
}
