package price_placements_feeds

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"strconv"
)

const (
	emptyFeed string = "feed is empty"
)

func statusCodeHandler(resp *http.Response) error {
	if resp == nil {
		return fmt.Errorf("не могу получить ответ сервера")
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("не могу получить ответ сервера. Statuscode: %v", resp.StatusCode)
	}
	return nil
}

type CustomInt64 struct {
	Int64 int64
	Valid bool
}

func (ci *CustomInt64) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var s string
	if err := d.DecodeElement(&s, &start); err != nil {
		return err
	}

	if s == "undefined" {
		ci.Valid = false
		return nil
	}

	customI, err := strconv.Atoi(s)
	if err != nil {
		return err
	}

	ci.Int64 = int64(customI)
	return nil
}
