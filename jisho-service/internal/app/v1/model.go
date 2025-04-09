package v1

type JishoResponse struct {
	Metadata JishoStatusResponse `json:"meta"`
	Data     []JishoData         `json:"data"`
}

type JishoStatusResponse struct {
	Status int `json:"status"`
}

type JishoData struct {
	Slug     string     `json:"slug"`
	IsCommon bool       `json:"is_common"`
	Tags     []string   `json:"tags"`
	JLPT     []string   `json:"jlpt"`
	Japanese []Japanese `json:"japanese"`
	Senses   []Sense    `json:"senses"`
}

type Japanese struct {
	Word    string `json:"word"`
	Reading string `json:"reading"`
}

type Sense struct {
	EnglishDefinitions []string `json:"english_definitions"`
	PartsOfSpeech      []string `json:"parts_of_speech"`
	Tags               []string `json:"tags"`
	SeeAlso            []string `json:"see_also"`
}
