package main

import "time"

// Status is like an enum and shows if a reviewable asset has been accepted, rejected or deleted
type Status int

const (
	// Preliminary means it has not been reviewed yet
	Preliminary Status = 1 + iota
	// Active means it has passed review
	Active
	// Outdated means it has once passed review, but is not relevant anymore
	Outdated
	// Deleted means a delete-request has been successful and reviewers thought this asset should be deleted
	Deleted
	// Rejected means it has not passed review
	Rejected
)

// Score is the sustainability rating of a scorable asset, either an "atomic" one (in a Rating), but usually, in a ScorableAsset, the averaged one
type Score struct {
	Environment int `json:"environment"` // range=[-100,100], air pollution, water pollution, ground pollution, waste, toxic substances released into environment etc., without GHG gases
	Climate     int `json:"climate"`     // range=[-100,100], emission of GHG gases and other climate-active actions like land-use change
	Society     int `json:"society"`     // range=[-100,100], working conditions, fair pay, workers' health, child labor, equity, treatment of suppliers, impact on society like charitable projects
	Health      int `json:"health"`      // range=[-100,100], impact on consumer's health, e.g. sugar and fat content in food or toxic substances in textiles or toys, acting on consumer
	Economy     int `json:"economy"`     // range=[-100,100], in the sense of 'value for money', longevity of product, price/performance ratio, is price too high because of the psychologically developed brand image? how economical is product for consumer?
}

// example:
// &Score{
//   Environment: -34,
//   Climate: -46,
//   Society: -7,
//   Health: -78,
//   Economy: 21,
// }

/**
  Reviewable assets must pass a peer review before going online. They can be
  edited/updated/deleted by other users, initiating another peer review.
**/

/* These are all abstract assets (see model), so don't need a DocType */

// ReviewableAsset is an asset that must pass a peer review before going online
type ReviewableAsset struct { /* only Comment extends ReviewableAsset */
	CreatedBy string    `json:"createdBy"` /* user ID, i.e. name */
	CreatedAt time.Time `json:"createdAt"`
	Status    Status    `json:"status"` // default=PRELIMINARY
}

// UpdatableAsset is an asset that must pass a peer review and can be changed later
type UpdatableAsset struct { /* only Information extends UpdatableAsset */
	ReviewableAsset
	UpdatedBy    string    `json:"updatedBy"` /* user ID, i.e. name */
	UpdatedAt    time.Time `json:"updatedAt"`
	Supersedes   string    `json:"supersedes"`   // optional /* (ID of) Previous version of this asset before it was updated. */
	SupersededBy string    `json:"supersededBy"` // optional /* (ID of) Newer version of this asset. */
	ChangeReason string    `json:"changeReason"` // optional
}

// ScorableAsset is an asset that must pass a review, can be changed later and holds a score (sustainability rating)
type ScorableAsset struct { /* all others (Product, Producer, Label) extend ScorableAsset */
	UpdatableAsset
	Score Score `json:"score"`
}
