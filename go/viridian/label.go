package viridian

// LabelLocaleData is the locale-specific (language-specific) part of a label
type LabelLocaleData struct {
	Lang        string   `json:"lang"` // regex=/^[a-z]{2}$/ // ISO language code according to https://en.wikipedia.org/wiki/ISO_639-1, there should be only one locale data for each language
	Name        string   `json:"name"`
	Description string   `json:"description"` // optional
	URL         string   `json:"url"`         // regex=/^[a-z]+:\/\/[^ ]+$/ optional
	Categories  []string `json:"categories"`
}

// Label is the asset representing a sustainability label that can be associated with a product or a producer, e.g. "Organic" or "Fairtrade"
type Label struct {
	ScorableAsset
	DocType string            `json:"docType"` // docType is used to distinguish the various types of objects in state database
	Locales []LabelLocaleData `json:"locales"`
	Version string            `json:"version"` // optional // label IDs
}
