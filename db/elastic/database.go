package elastic

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/olivere/elastic"

	"github.com/alerting/go-cap"
	"github.com/alerting/go-cap-process/db"
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

func (es *Elastic) AddAlert(alerts ...*cap.Alert) error {
	bulkAlert := es.client.Bulk().Index(es.index).Type("_doc")
	bulkInfo := es.client.Bulk().Index(es.index).Type("_doc")

	for _, alert := range alerts {
		// Convert to map[string]interface{}
		var alertMap map[string]interface{}
		b, _ := json.Marshal(&alert)
		json.Unmarshal(b, &alertMap)

		// We don't need the infos item (will be added independently)
		delete(alertMap, "infos")

		// Setup Parent
		alertMap["_object"] = map[string]string{
			"name": "alert",
		}

		bulkAlert.Add(elastic.NewBulkIndexRequest().Id(alert.Id()).Doc(alertMap))

		for indx, info := range alert.Infos {
			var infoMap map[string]interface{}
			b, _ := json.Marshal(&info)
			json.Unmarshal(b, &infoMap)

			// Setup Parent
			infoMap["_object"] = map[string]string{
				"name":   "info",
				"parent": alert.Id(),
			}

			bulkInfo.Add(
				elastic.NewBulkIndexRequest().
					Id(fmt.Sprintf("%s:%d", alert.Id(), indx)).
					Routing(alert.Id()).
					Doc(infoMap))
		}
	}

	// TODO: Process errors
	_, err := bulkAlert.Do(context.Background())
	if err != nil {
		return err
	}

	res, err := bulkInfo.Do(context.Background())
	if err != nil {
		return err
	}
	if res.Errors {
		for _, i := range res.Items {
			if i["index"].Error != nil {
				fmt.Println(i["index"].Id)
				fmt.Println(i["index"].Error)
			}
		}
	}

	return nil
}

func (es *Elastic) AlertExists(reference *cap.Reference) (bool, error) {
	item := elastic.NewMultiGetItem().Index(es.index).Type("_doc").Id(reference.Id())
	res, err := es.client.MultiGet().Add(item).Do(context.Background())
	if err != nil {
		return false, err
	}

	return res.Docs[0].Found, nil
}

func (es *Elastic) GetAlert(reference *cap.Reference) (*cap.Alert, error) {
	return es.GetAlertById(reference.Id())
}

func (es *Elastic) GetAlertById(id string) (*cap.Alert, error) {
	item, err := es.client.Get().Index(es.index).Type("_doc").Id(id).Do(context.Background())
	if err != nil {
		return nil, err
	}

	// Fetch the alert itself
	var alert cap.Alert
	err = json.Unmarshal(*item.Source, &alert)
	if err != nil {
		return nil, err
	}

	// Fetch the children (ie. infos)
	finder := es.NewInfoFinder()
	finder = finder.AlertId(alert.Id())
	finder = finder.Sort("_id")

	infos, err := finder.Find()
	if err != nil {
		return nil, err
	}

	for _, hit := range infos.Hits {
		alert.Infos = append(alert.Infos, *hit.Info)
	}

	return &alert, nil
}

func (es *Elastic) NewInfoFinder() db.InfoFinder {
	return NewInfoFinder(es)
}
