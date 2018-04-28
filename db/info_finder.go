package db

import cap "github.com/alerting/go-cap"
import "time"

type InfoHit struct {
	Id      string    `json:"id"`
	AlertId string    `json:"alert_id"`
	Info    *cap.Info `json:"info"`
}

type InfoResults struct {
	TotalHits int64      `json:"total_hits"`
	Hits      []*InfoHit `json:"hits"`
}

type InfoFinder interface {
	// Filter
	Status(status cap.Status) InfoFinder
	MessageType(messageType cap.MessageType) InfoFinder
	Scope(scope cap.Scope) InfoFinder

	Language(language string) InfoFinder
	Certainty(certainty cap.Certainty) InfoFinder
	Severity(severity cap.Severity) InfoFinder
	Urgency(urgency cap.Urgency) InfoFinder
	Headline(headline string) InfoFinder
	Description(description string) InfoFinder
	Instruction(instruction string) InfoFinder
	EffectiveGte(t time.Time) InfoFinder
	EffectiveGt(t time.Time) InfoFinder
	EffectiveLte(t time.Time) InfoFinder
	EffectiveLt(t time.Time) InfoFinder
	ExpiresGte(t time.Time) InfoFinder
	ExpiresGt(t time.Time) InfoFinder
	ExpiresLte(t time.Time) InfoFinder
	ExpiresLt(t time.Time) InfoFinder
	OnsetGte(t time.Time) InfoFinder
	OnsetGt(t time.Time) InfoFinder
	OnsetLte(t time.Time) InfoFinder
	OnsetLt(t time.Time) InfoFinder

	Area(area string) InfoFinder
	Point(lat, lon float64) InfoFinder

	// Pagination
	Start(start int) InfoFinder
	Count(count int) InfoFinder

	// Sorting
	Sort(fields ...string) InfoFinder

	Find() (*InfoResults, error)
}
