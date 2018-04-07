package db

import cap "github.com/alerting/go-cap"

type InfoHit struct {
	AlertId string    `json:"alert_id"`
	Info    *cap.Info `json:"info"`
}

type InfoResults struct {
	TotalHits int64      `json:"total_hits"`
	Hits      []*InfoHit `json:"hits"`
}

type InfoFinder interface {
	// Filter
	Language(language string) InfoFinder
	Certainty(certainty cap.Certainty) InfoFinder
	Severity(severity cap.Severity) InfoFinder
	Urgency(urgency cap.Urgency) InfoFinder
	Headline(headline string) InfoFinder
	Description(description string) InfoFinder
	Instruction(instruction string) InfoFinder

	// Pagination
	Start(start int) InfoFinder
	Count(count int) InfoFinder

	// Sorting
	Sort(fields ...string) InfoFinder

	Find() (*InfoResults, error)
}
