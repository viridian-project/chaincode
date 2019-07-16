package main

type Producer struct {
	ScorableAsset
	DocType string   `json:"docType"` // docType is used to distinguish the various types of objects in state database
	Name    string   `json:"name"`
	Address string   `json:"address"` // optional
	Url     string   `json:"url"`     // optional
	Labels  []string `json:"labels"`
}
