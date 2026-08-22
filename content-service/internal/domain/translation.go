package domain

import "time"

type TranslationEntryInput struct {
	OriginalHTML    string    `json:"originalHtml"`
	TranslationHTML string    `json:"translationHtml"`
	CreatedAt       time.Time `json:"createdAt"`
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

type TranslationEntryOutput struct {
	ID              int32     `json:"id"`
	OriginalHTML    string    `json:"originalHtml"`
	TranslationHTML string    `json:"translationHtml"`
	CreatedAt       time.Time `json:"createdAt"`
}

type TranslationListOutput struct {
	ID        int32                    `json:"id"`
	Name      string                   `json:"name"`
	CreatedAt time.Time                `json:"createdAt"`
	Entries   []TranslationEntryOutput `json:"entries"`
}

type GetListsResult struct {
	Lists []TranslationListOutput `json:"lists"`
}

type DeleteListRequest struct {
	ID int32 `json:"id"`
}

type DeleteListResult struct {
	ID int32 `json:"id"`
}
