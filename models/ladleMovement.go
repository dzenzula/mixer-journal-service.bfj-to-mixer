package models

type LadleMovement struct {
	Date           string `json:"date"`
	NumDp          int    `json:"numDp"`
	NumTapping     int    `json:"numTapping"`
	LadleTapping   string `json:"ladleTapping"`
	DtCloseTaphole string `json:"dtCloseTaphole"`
	//DtPlum            string `json:"dtPlum"`
	MassCastIron      float64 `json:"massCastIron"`
	TemperExhaustIron int     `json:"temperExhaustIron"`
	//DtEndDrainingDc   string `json:"dtEndDrainingDc"`
}
