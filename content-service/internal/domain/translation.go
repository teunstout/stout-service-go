package domain

type TranslationEntryInput struct {
	OriginalHTML    string `json:"originalHtml"`
	TranslationHTML string `json:"translationHtml"`
}

type SyncListRequest struct {
	ID      *int32                  `json:"id"`
	Name    string                  `json:"name"`
	Entries []TranslationEntryInput `json:"entries"`
}

type SyncListResult struct {
	ID         int32  `json:"id"`
	Name       string `json:"name"`
	EntryCount int    `json:"entryCount"`
}
