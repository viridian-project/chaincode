package main

type InfoCategory int const (
	GENERAL_INFORMATION InfoCategory = 1 + iota
	LIFE_CYCLE_ANALYSIS
	EXTERNAL_COSTS
	STUDY_OR_PAPER
	PRESS_ARTICLE
	INVESTIGATIVE_REPORT
	CORPORATE_SOCIAL_RESPONSIBILITY
	JURISDICTION
	OTHER
)
type Information struct {
	UpdatableAsset
	Title string `json:"title"`
	Category InfoCategory `json:"category"`
	Description string `json:"description"`
	Sources []Source `json:"sources"`
	Weight int32 `json:"weight"`
}
