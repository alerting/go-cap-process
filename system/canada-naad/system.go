package canadanaad

import (
	"encoding/xml"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/alerting/go-cap"
)

type CanadaNAAD struct {
	Fetch *url.URL
}

func CreateSystem(fetch string) (*CanadaNAAD, error) {
	fetchUrl, err := url.Parse(fetch)
	if err != nil {
		return nil, err
	}

	return &CanadaNAAD{
		Fetch: fetchUrl,
	}, nil
}

func clean(str string) string {
	str = strings.Replace(str, "-", "_", -1)
	str = strings.Replace(str, "+", "p", -1)
	str = strings.Replace(str, ":", "_", -1)

	return str
}

func (sys *CanadaNAAD) GetAlert(reference *cap.Reference) (*cap.Alert, error) {
	// Generate the URL
	u, err := url.Parse(fmt.Sprintf("%s/%sI%s.xml",
		reference.Sent.Format("2006-01-02"),
		clean(reference.Sent.FormatCAP()),
		clean(reference.Identifier)))

	if err != nil {
		return nil, err
	}

	u = sys.Fetch.ResolveReference(u)

	res, err := http.Get(u.String())
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, errors.New(fmt.Sprintf("Expected status 200, got %s", res.Status))
	}

	d := xml.NewDecoder(res.Body)

	var alert cap.Alert
	if err = d.Decode(&alert); err != nil {
		return nil, err
	}

	return &alert, nil
}
