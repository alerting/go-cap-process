package elastic

import (
	"context"
	"encoding/json"
	"github.com/alerting/go-cap"
	"github.com/alerting/go-cap-process/db"
	"github.com/olivere/elastic"
	"strings"
	//"log"
)

type InfoFinder struct {
	elastic *Elastic

	parentFields map[string]string
	termFields   map[string]string
	textFields   map[string]string
	area         string
	point        *elastic.GeoPoint

	start int
	count int

	sort []string
}

func NewInfoFinder(elastic *Elastic) db.InfoFinder {
	return &InfoFinder{
		elastic:      elastic,
		parentFields: make(map[string]string),
		termFields:   make(map[string]string),
		textFields:   make(map[string]string),
		start:        -1,
		count:        -1,
		sort:         make([]string, 0),
	}
}

/** FILTERS **/
func (f *InfoFinder) Status(status cap.Status) db.InfoFinder {
	f.parentFields["status"] = status.String()
	return f
}

func (f *InfoFinder) MessageType(messageType cap.MessageType) db.InfoFinder {
	f.parentFields["message_type"] = messageType.String()
	return f
}

func (f *InfoFinder) Scope(scope cap.Scope) db.InfoFinder {
	f.parentFields["scope"] = scope.String()
	return f
}

func (f *InfoFinder) Language(language string) db.InfoFinder {
	f.termFields["language"] = language
	return f
}

func (f *InfoFinder) Certainty(certainty cap.Certainty) db.InfoFinder {
	f.termFields["certainty"] = certainty.String()
	return f
}

func (f *InfoFinder) Severity(severity cap.Severity) db.InfoFinder {
	f.termFields["severity"] = severity.String()
	return f
}

func (f *InfoFinder) Urgency(urgency cap.Urgency) db.InfoFinder {
	f.termFields["urgency"] = urgency.String()
	return f
}

func (f *InfoFinder) Headline(headline string) db.InfoFinder {
	f.textFields["headline"] = headline
	return f
}

func (f *InfoFinder) Description(description string) db.InfoFinder {
	f.textFields["description"] = description
	return f
}

func (f *InfoFinder) Instruction(instruction string) db.InfoFinder {
	f.textFields["instruction"] = instruction
	return f
}

func (f *InfoFinder) Area(area string) db.InfoFinder {
	f.area = area
	return f
}

func (f *InfoFinder) Point(lat, lon float64) db.InfoFinder {
	f.point = elastic.GeoPointFromLatLon(lat, lon)
	return f
}

/** PAGINATION **/
func (f *InfoFinder) Start(start int) db.InfoFinder {
	f.start = start
	return f
}

func (f *InfoFinder) Count(count int) db.InfoFinder {
	f.count = count
	return f
}

/** SORTING **/
func (f *InfoFinder) Sort(fields ...string) db.InfoFinder {
	f.sort = append(f.sort, fields...)
	return f
}

/** FIND **/
func (f *InfoFinder) Find() (*db.InfoResults, error) {
	search := f.elastic.client.Search().Index(f.elastic.index).Type("_doc")
	search = f.query(search)
	search = f.pagination(search)
	search = f.sorting(search)

	res, err := search.Do(context.Background())
	if err != nil {
		return nil, err
	}

	// Process results
	results := db.InfoResults{
		TotalHits: res.Hits.TotalHits,
		Hits:      make([]*db.InfoHit, 0),
	}

	for _, hit := range res.Hits.Hits {
		var info cap.Info
		if err = json.Unmarshal(*hit.Source, &info); err != nil {
			return nil, err
		}

		infoHit := db.InfoHit{
			Id:      hit.Id,
			AlertId: hit.Routing,
			Info:    &info,
		}

		results.Hits = append(results.Hits, &infoHit)
	}

	return &results, nil
}

func (f *InfoFinder) query(service *elastic.SearchService) *elastic.SearchService {
	q := elastic.NewBoolQuery()

	// Parent filter
	if len(f.parentFields) > 0 {
		pq := elastic.NewBoolQuery()

		for k, v := range f.parentFields {
			pq = pq.Must(elastic.NewTermQuery(k, v))
		}

		q = q.Must(elastic.NewHasParentQuery("alert", pq))
	} else {
		q = q.Must(elastic.NewHasParentQuery("alert", elastic.NewMatchAllQuery()))
	}

	// Filter on termFields
	if len(f.termFields) > 0 {
		for k, v := range f.termFields {
			q = q.Must(elastic.NewTermQuery(k, v))
		}
	}

	// Filter on textFields
	if len(f.textFields) > 0 {
		for k, v := range f.textFields {
			q = q.Must(elastic.NewQueryStringQuery(v).Field(k))
		}
	}

	// Filter on area
	if f.area != "" || f.point != nil {
		aq := elastic.NewBoolQuery()

		if f.area != "" {
			aq = aq.Must(elastic.NewQueryStringQuery(f.area).Field("areas.description"))
		}

		if f.point != nil {
			pq := elastic.NewBoolQuery()
			pq = pq.Should(NewGeoShapeQuery("areas.polygons").SetPoint(f.point.Lat, f.point.Lon))
			pq = pq.Should(NewGeoShapeQuery("areas.circles").SetPoint(f.point.Lat, f.point.Lon))

			aq = aq.Must(pq)
		}

		nq := elastic.NewNestedQuery("areas", aq)
		nq.InnerHit(elastic.NewInnerHit().FetchSourceContext(elastic.NewFetchSourceContext(false)))

		q = q.Must(nq)
	}

	service = service.Query(q)
	return service
}

func (f *InfoFinder) pagination(service *elastic.SearchService) *elastic.SearchService {
	if f.start >= 0 {
		service = service.From(f.start)
	}

	if f.count >= 0 {
		service = service.Size(f.count)
	}

	return service
}

func (f *InfoFinder) sorting(service *elastic.SearchService) *elastic.SearchService {
	if len(f.sort) == 0 {
		service = service.Sort("_score", false)
		return service
	}

	// Prefix of "-" means to sort descending.
	for _, field := range f.sort {
		asc := true
		if strings.HasPrefix(field, "-") {
			field = field[1:]
			asc = false
		}

		service = service.Sort(field, asc)
	}

	return service
}
