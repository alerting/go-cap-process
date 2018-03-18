package elastic

import (
	"context"
	"errors"

	"github.com/olivere/elastic"

	"github.com/alerting/go-cap"
)

type Elastic struct {
	client *elastic.Client
	index  string
}

func CreateDatabase(url string, index string) (*Elastic, error) {
	db := Elastic{
		index: index,
	}

	var err error
	db.client, err = elastic.NewClient(elastic.SetURL(url))
	if err != nil {
		return nil, err
	}

	return &db, nil
}

func (es *Elastic) Setup() error {
	exists, err := es.client.IndexExists(es.index).Do(context.Background())
	if err != nil {
		return err
	}

	if !exists {
		_, err = es.client.CreateIndex(es.index).BodyString(mapping).Do(context.Background())
		if err != nil {
			return err
		}
	}

	return nil
}

func (es *Elastic) AddAlert(alert *cap.Alert) error {
	_, err := es.client.Index().
		Index(es.index).Type("alert").Id(alert.Id()).
		BodyJson(alert).
		Refresh("wait_for").
		Do(context.Background())

	return err
}

func (es *Elastic) AlertExists(reference *cap.Reference) (bool, error) {
	item := elastic.NewMultiGetItem().Index(es.index).Type("alert").Id(reference.Id())
	res, err := es.client.MultiGet().Add(item).Do(context.Background())
	if err != nil {
		return false, err
	}

	return res.Docs[0].Found, nil
}

func (es *Elastic) GetAlert(reference *cap.Reference) (*cap.Alert, error) {
	return nil, errors.New("Not implemented")
}
