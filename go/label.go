package main

type LabelLocaleData struct {
	Lang        string   `json:"lang"` // regex=/^[a-z]{2}$/ // ISO language code according to https://en.wikipedia.org/wiki/ISO_639-1, there should be only one locale data for each language
	Name        string   `json:"name"`
	Description string   `json:"description"` // optional
	Url         string   `json:"url"`         // regex=/^[a-z]+:\/\/[^ ]+$/ optional
	Categories  []string `json:"categories"`
}

type Label struct {
	ScorableAsset
	DocType string            `json:"docType"` // docType is used to distinguish the various types of objects in state database
	Locales []LabelLocaleData `json:"locales"`
	Version string            `json:"version"` // optional // label IDs
}
