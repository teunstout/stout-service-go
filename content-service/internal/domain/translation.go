package domain

import "time"

type TranslationEntryInput struct {
	ID              *int32    `json:"id"`
	OriginalHTML    string    `json:"originalHtml"`
	TranslationHTML string    `json:"translationHtml"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`
}

type SyncListRequest struct {
	ID      *int32                  `json:"id"`
	Name    string                  `json:"name"`
	Entries []TranslationEntryInput `json:"entries"`
}

type SyncEntryResult struct {
	ID int32 `json:"id"`
}

type SyncListResult struct {
	ID      int32             `json:"id"`
	Name    string            `json:"name"`
	Entries []SyncEntryResult `json:"entries"`
}

type TranslationEntryOutput struct {
	ID              int32     `json:"id"`
	OriginalHTML    string    `json:"originalHtml"`
	TranslationHTML string    `json:"translationHtml"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`
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

type DeleteEntriesRequest struct {
	IDs []int32 `json:"ids"`
}

type DeleteEntriesResult struct {
	DeletedIDs []int32 `json:"deletedIds"`
}
