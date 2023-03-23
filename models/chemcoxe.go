package models

import "time"

type ChemCoxe struct {
	ID              int       `json:"id"`
	IDJournal       int       `json:"idJournal"`
	SupplierID      *int      `json:"supplierId"`
	Supplier        string    `json:"suplier"`
	TechAnHumidity  float64   `json:"techAnHumidity"`
	TechAnAsh       float64   `json:"techAnAsh"`
	TechAnVolatiles float64   `json:"techAnVolatiles"`
	TechAnSulphur   float64   `json:"techAnSulphur"`
	StrengthM25     float64   `json:"strengthM25"`
	StrengthM10     float64   `json:"strengthM10"`
	StrengthGp      float64   `json:"strengthGp"`
	StrengthPc      float64   `json:"strengthPc"`
	CompositionA    float64   `json:"compositionA"`
	CompositionB    float64   `json:"compositionB"`
	CompositionC    float64   `json:"compositionC"`
	CompositionD    float64   `json:"compositionD"`
	CompositionE    float64   `json:"compositionE"`
	DtUnload        time.Time `json:"dtUnload"`
	AmountP         int       `json:"amountP"`
	SupplierUnload  string    `json:"suplierUnload"`
	DtInsert        time.Time `json:"dtInsert"`
	DtUpdate        time.Time `json:"dtUpdate"`
	CreatedOn       time.Time `json:"createdOn"`
	ModifyBy        time.Time `json:"modifyBy"`
}
