package main

// Producer is the asset associated with bringing a product to market, so being responsible for it
type Producer struct {
	ScorableAsset
	DocType string   `json:"docType"` // docType is used to distinguish the various types of objects in state database
	Name    string   `json:"name"`
	Address string   `json:"address"` // optional
	URL     string   `json:"url"`     // optional
	Labels  []string `json:"labels"`
}
