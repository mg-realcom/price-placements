package price_placements_feeds

import (
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"strconv"
	"strings"
	"time"
)

type CianFeed struct {
	LastModified time.Time
	FeedVersion  string   `xml:"feed_version"`
	Object       []Object `xml:"object"`
}

type Object struct {
	ExternalId  string `xml:"ExternalId"`
	Description string `xml:"Description"`
	Address     string `xml:"Address"`
	Coordinates struct {
		Lat float32 `xml:"Lat"`
		Lng float32 `xml:"Lng"`
	} `xml:"Coordinates"`
	CadastralNumber string `xml:"CadastralNumber"`
	Phones          struct {
		PhoneSchema struct {
			CountryCode string `xml:"CountryCode"`
			Number      string `xml:"Number"`
		} `xml:"PhoneSchema"`
	} `xml:"Phones"`
	LayoutPhoto struct {
		IsDefault bool   `xml:"IsDefault"`
		FullUrl   string `xml:"FullUrl"`
	} `xml:"LayoutPhoto"`
	Photos struct {
		PhotoSchema []struct {
			FullUrl   string `xml:"FullUrl"`
			IsDefault bool   `xml:"IsDefault"`
		} `xml:"PhotoSchema"`
	} `xml:"Photos"`
	Category              string  `xml:"Category"`
	RoomType              string  `xml:"RoomType"`
	FlatRoomsCount        int64   `xml:"FlatRoomsCount"`
	TotalArea             float32 `xml:"TotalArea"`
	LivingArea            float32 `xml:"LivingArea"`
	KitchenArea           float32 `xml:"KitchenArea"`
	ProjectDeclarationUrl string  `xml:"ProjectDeclarationUrl"`
	FloorNumber           int64   `xml:"FloorNumber"`
	CombinedWcsCount      int64   `xml:"CombinedWcsCount"`
	Building              struct {
		FloorsCount         int64  `xml:"FloorsCount"`
		MaterialType        string `xml:"MaterialType"`
		PassengerLiftsCount int64  `xml:"PassengerLiftsCount"`
		CargoLiftsCount     int64  `xml:"CargoLiftsCount"`
		Parking             struct {
			Type string `xml:"Type"`
		} `xml:"Parking"`
		Deadline struct {
			Quarter    string `xml:"Quarter"`
			Year       int64  `xml:"Year"`
			IsComplete bool   `xml:"IsComplete"`
		} `xml:"Deadline"`
	} `xml:"Building"`
	BargainTerms struct {
		Price           CustomFloat64 `xml:"Price"`
		Currency        string        `xml:"Currency"`
		MortgageAllowed bool          `xml:"MortgageAllowed"`
		SaleType        string        `xml:"SaleType"`
	} `xml:"BargainTerms"`
	JKSchema struct {
		ID    int32  `xml:"Id"`
		Name  string `xml:"Name"`
		House struct {
			ID   int32  `xml:"Id"`
			Name string `xml:"Name"`
			Flat struct {
				FlatNumber    int32  `xml:"FlatNumber"`
				SectionNumber string `xml:"SectionNumber"`
				FlatType      string `xml:"FlatType"`
			} `xml:"Flat"`
		} `xml:"House"`
	} `xml:"JKSchema"`
	Decoration      string  `xml:"Decoration"`
	WindowsViewType string  `xml:"WindowsViewType"`
	CeilingHeight   float32 `xml:"CeilingHeight"`
	Undergrounds    struct {
		UndergroundInfoSchema []struct {
			TransportType string `xml:"TransportType"`
			Time          int64  `xml:"Time"`
			ID            int64  `xml:"Id"`
		} `xml:"UndergroundInfoSchema"`
	} `xml:"Undergrounds"`
	IsApartments bool `xml:"isApartments"`
}

type CustomFloat64 struct {
	Float64 float64
}

func (cf *CustomFloat64) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var s string
	if err := d.DecodeElement(&s, &start); err != nil {
		return err
	}
	s = strings.ReplaceAll(s, ",", ".")
	float, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return err
	}
	cf.Float64 = float
	return nil
}

func (f *CianFeed) Get(url string) (err error) {
	resp, err := GetResponse(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	err = statusCodeHandler(resp)
	if err != nil {
		return err
	}

	AttributeLastModified := resp.Header.Get("Last-Modified")
	if AttributeLastModified != "" {
		lastModifiedDate, err := time.Parse(time.RFC1123, resp.Header.Get("Last-Modified"))
		if err != nil {
			return err
		}
		f.LastModified = lastModifiedDate
	} else {
		log.Println("Header not contains `Last-Modified`")
	}
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	err = xml.Unmarshal(responseBody, &f)
	if err != nil {
		return err
	}
	//TODO Исправить значение f.LastModified
	return nil
}

func (f *CianFeed) Check() (results []string) {
	if len(f.Object) == 0 {
		results = append(results, emptyFeed)
		return results
	}

	if len(f.Object) <= 10 {
		results = append(results, fmt.Sprintf("feed contains only %v items", len(f.Object)))
		return results
	}
	for idx, lot := range f.Object {
		if lot.ExternalId == "" {
			results = append(results, fmt.Sprintf("field ExternalId is empty. Position: %v", idx))
		}
		if lot.Category == "" {
			results = append(results, fmt.Sprintf("field Category is empty. Position: %v", idx))
		}
		if lot.FlatRoomsCount == 0 {
			results = append(results, fmt.Sprintf("field FlatRoomsCount is empty. Position: %v", idx))
		}
		if lot.TotalArea == 0 {
			results = append(results, fmt.Sprintf("field TotalArea is empty. Position: %v", idx))
		}
		if lot.FloorNumber == 0 {
			results = append(results, fmt.Sprintf("field FloorNumber is empty. Position: %v", idx))
		}
		if lot.Building.FloorsCount == 0 {
			results = append(results, fmt.Sprintf("field Building.FloorsCount is empty. Position: %v", idx))
		}
		if lot.JKSchema.ID == 0 {
			results = append(results, fmt.Sprintf("field JKSchema.ID is empty. Position: %v", idx))
		}
		if lot.JKSchema.Name == "" {
			results = append(results, fmt.Sprintf("field JKSchema.Name is empty. Position: %v", idx))
		}
		if lot.JKSchema.House.ID == 0 {
			results = append(results, fmt.Sprintf("field JKSchema.House.ID is empty. Position: %v", idx))
		}
		if lot.JKSchema.House.Name == "" {
			results = append(results, fmt.Sprintf("field JKSchema.House.Name is empty. Position: %v", idx))
		}
		if lot.BargainTerms.Price.Float64 == 0 {
			results = append(results, fmt.Sprintf("field BargainTerms.Price is empty. Position: %v", idx))
		}
	}
	return results
}
