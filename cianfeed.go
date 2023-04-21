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
		PhotoSchema []PhotoSchema `xml:"PhotoSchema"`
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

type PhotoSchema struct {
	FullUrl   string `xml:"FullUrl"`
	IsDefault bool   `xml:"IsDefault"`
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
	if len(f.Object) < 2 {
		results = append(results, emptyFeed)
		return results
	}

	if len(f.Object) <= 10 {
		results = append(results, fmt.Sprintf("feed contains only %v items", len(f.Object)))
		return results
	}
	for idx, lot := range f.Object {
		id := lot.ExternalId

		if lot.ExternalId == "" {
			results = append(results, fmt.Sprintf("field ExternalId is empty. Position: %v", idx))
		}
		checkStringWithID(id, "object", "Address", lot.Address, &results)
		checkStringWithID(id, "object.Phones.PhoneSchema", "CountryCode", lot.Phones.PhoneSchema.CountryCode, &results)
		checkStringWithID(id, "object.Phones.PhoneSchema", "Number", lot.Phones.PhoneSchema.Number, &results)
		checkStringWithID(id, "object.LayoutPhoto.FullUrl", "IsDefault", lot.LayoutPhoto.FullUrl, &results)
		checkStringWithID(id, "object", "Category", lot.Category, &results)

		for idx, photoSchema := range lot.Photos.PhotoSchema {
			checkStringWithPos(idx, "object.Photos.PhotoSchema", "FullUrl", photoSchema.FullUrl, &results)
		}

		checkZeroWithID(id, "object", "FlatRoomsCount", int(lot.FlatRoomsCount), &results)
		checkZeroWithID(id, "object", "TotalArea", int(lot.TotalArea), &results)
		checkZeroWithID(id, "object", "FloorNumber", int(lot.FloorNumber), &results)
		checkZeroWithID(id, "object.Building", "FloorsCount", int(lot.Building.FloorsCount), &results)
		checkZeroWithID(id, "object.Building.Deadline", "Year", int(lot.Building.Deadline.Year), &results)
		checkStringWithID(id, "object.Building.Deadline", "Quarter", lot.Building.Deadline.Quarter, &results)
		checkZeroWithID(id, "object.BargainTerms.Price", "Price", int(lot.BargainTerms.Price.Float64), &results)
		checkZeroWithID(id, "object.JKSchema", "Id", int(lot.JKSchema.ID), &results)
		checkStringWithID(id, "object.JKSchema", "Name", lot.JKSchema.Name, &results)
		checkZeroWithID(id, "object.JKSchema.House", "Id", int(lot.JKSchema.House.ID), &results)
		checkStringWithID(id, "object.JKSchema.House", "Name", lot.JKSchema.House.Name, &results)

		if lot.Building.Deadline.Year < int64(time.Now().Year()) && lot.Building.Deadline.IsComplete == false {
			results = append(results, fmt.Sprintf("field Building.Deadline is False for %v. InternalID: %v", lot.Building.Deadline.Year, lot.ExternalId))
		}
		if lot.FloorNumber > lot.Building.FloorsCount {
			results = append(results, fmt.Sprintf("field FloorNumber is greater than Building.FloorsCount. InternalID: %v", lot.ExternalId))
		}
		if len(lot.Photos.PhotoSchema) < 3 {
			results = append(results, fmt.Sprintf("field Photos.PhotoSchema contains '%v' items. InternalID: %v", len(lot.Photos.PhotoSchema), lot.ExternalId))
		}
	}

	return results
}
