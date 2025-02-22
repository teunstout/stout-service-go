package v1

type JishoResponse struct {
	Metadata JishoStatusResponse `json:"meta"`
	Data     []JishoData         `json:"data"`
}

type JishoStatusResponse struct {
	Status int `json:"status"`
}

type JishoData struct {
	Slug        string      `json:"slug"`
	IsCommon    bool        `json:"is_common"`
	Tags        []string    `json:"tags"`
	JLPT        []string    `json:"jlpt"`
	Japanese    []Japanese  `json:"japanese"`
	Senses      []Sense     `json:"senses"`
	Attribution Attribution `json:"attribution"`
}

type Japanese struct {
	Word    string `json:"word,omitempty"`
	Reading string `json:"reading"`
}

type Sense struct {
	EnglishDefinitions []string `json:"english_definitions"`
	PartsOfSpeech      []string `json:"parts_of_speech"`
	Links              []string `json:"links"`
	Tags               []string `json:"tags"`
	Restrictions       []string `json:"restrictions"`
	SeeAlso            []string `json:"see_also"`
	Antonyms           []string `json:"antonyms"`
	Source             []string `json:"source"`
	Info               []string `json:"info"`
}

type Attribution struct {
	JMDict   bool `json:"jmdict"`
	JMNedict bool `json:"jmnedict"`
	DBPedia  bool `json:"dbpedia"`
}
