package resolvedschema

// Item -
type Item struct {
	Nome       string `json:"nome,omitempty"`
	Tipo       string `json:"tipo,omitempty"`
	Quantidade int    `json:"quantidade,omitempty"`
}

// Test -
type Test struct {
	Nome  string `json:"nome,omitempty"`
	Test1 Test1  `json:"test_1,omitempty"`
}

// Test1 -
type Test1 struct {
	ID   string `json:"id,omitempty"`
	Alvo string `json:"alvo,omitempty"`
}
