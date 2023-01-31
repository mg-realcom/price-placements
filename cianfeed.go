package price_placements_feeds

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"log"
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
		Price           float64 `xml:"Price"`
		Currency        string  `xml:"Currency"`
		MortgageAllowed bool    `xml:"MortgageAllowed"`
		SaleType        string  `xml:"SaleType"`
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

func (f *CianFeed) Get(url string) (err error) {
	resp, err := GetResponse(url)
	defer resp.Body.Close()
	if err != nil {
		return err
	}

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

func (f *CianFeed) Check() (errs []error) {
	if len(f.Object) == 0 {
		errs = append(errs, errors.New("feed is empty"))
		return errs
	}
	for idx, lot := range f.Object {
		if lot.ExternalId == "" {
			errs = append(errs, fmt.Errorf("field ExternalId is empty. Position: %v", idx))
		}
		if lot.Category == "" {
			errs = append(errs, fmt.Errorf("field Category is empty. Position: %v", idx))
		}
		if lot.FlatRoomsCount == 0 {
			errs = append(errs, fmt.Errorf("field FlatRoomsCount is empty. Position: %v", idx))
		}
		if lot.TotalArea == 0 {
			errs = append(errs, fmt.Errorf("field TotalArea is empty. Position: %v", idx))
		}
		if lot.FloorNumber == 0 {
			errs = append(errs, fmt.Errorf("field FloorNumber is empty. Position: %v", idx))
		}
		if lot.Building.FloorsCount == 0 {
			errs = append(errs, fmt.Errorf("field Building.FloorsCount is empty. Position: %v", idx))
		}
		if lot.JKSchema.ID == 0 {
			errs = append(errs, fmt.Errorf("field JKSchema.ID is empty. Position: %v", idx))
		}
		if lot.JKSchema.Name == "" {
			errs = append(errs, fmt.Errorf("field JKSchema.Name is empty. Position: %v", idx))
		}
		if lot.JKSchema.House.ID == 0 {
			errs = append(errs, fmt.Errorf("field JKSchema.House.ID is empty. Position: %v", idx))
		}
		if lot.JKSchema.House.Name == "" {
			errs = append(errs, fmt.Errorf("field JKSchema.House.Name is empty. Position: %v", idx))
		}
		if lot.BargainTerms.Price == 0 {
			errs = append(errs, fmt.Errorf("field BargainTerms.Price is empty. Position: %v", idx))
		}
	}
	return errs
}
