package cpanmeta

type Prereqs struct {
	Configure Phase `json:"configure"`
	Runtime   Phase `json:"runtime"`
	Build     Phase `json:"build"`
	Test      Phase `json:"test"`
	Develop   Phase `json:"develop"`
}
