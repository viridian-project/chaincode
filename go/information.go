package main

import "time"

type InfoCategory int

const (
	GeneralInformation InfoCategory = 1 + iota
	LifeCycleAnalysis
	ExternalCosts
	StudyOrPaper
	PressArticle
	InvestigativeReport
	CorporateSocialResponsibility
	Jurisdiction
	Other
)

type Source struct {
}
type WebSource struct {
	Source
	URL        string    `json:"url"` // regex=/^[a-z]+:\/\/[^ ]+$/
	AccessDate time.Time `json:"accessDate"`
	Title      string    `json:"title"`   // optional
	Authors    []string  `json:"authors"` // optional
}
type BookSource struct {
	Source
	Title       string   `json:"title"`
	Authors     []string `json:"authors"`
	PublishYear int      `json:"publishYear"`
	Publisher   string   `json:"publisher"` // optional
	ISBN        string   `json:"authors"`   // optional
	Pages       []int    `json:"authors"`   // range=[1,] optional
	URL         string   `json:"authors"`   // regex=/^[a-z]+:\/\/[^ ]+$/ optional
}
type ArticleSource struct {
	Source
	Title     string   `json:"title"`
	Authors   []string `json:"authors"`
	Journal   string   `json:"journal"`
	Year      int      `json:"year"`
	Month     int      `json:"month"`     // range=[1,12] optional
	Volume    int      `json:"volume"`    // optional
	FirstPage int      `json:"firstPage"` // optional
	LastPage  int      `json:"lastPage"`  // optional
	DOI       string   `json:"doi"`       // optional
	URL       string   `json:"url"`       // regex=/^[a-z]+:\/\/[^ ]+$/ optional
	BookTitle string   `json:"bookTitle"` // optional
	Editor    string   `json:"editor"`    // optional
}

type Information struct {
	UpdatableAsset
	Title       string       `json:"title"`
	Category    InfoCategory `json:"category"`
	Target      string       `json:"target"` // ID of the targeted scorable asset
	Description string       `json:"description"`
	Sources     []Source     `json:"sources"`
	Weight      int32        `json:"weight"`
}
