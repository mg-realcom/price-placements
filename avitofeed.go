package price_placements_feeds

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"log"
	"time"
)

type AvitoFeed struct {
	LastModified  time.Time
	XMLName       xml.Name `xml:"Ads"`
	FormatVersion int      `xml:"formatVersion,attr"`
	Target        string   `xml:"target,attr"`
	Ad            []Ad     `xml:"Ad"`
}

type Ad struct {
	ID              string  `xml:"Id"`
	AdStatus        string  `xml:"AdStatus"`
	AllowEmail      string  `xml:"AllowEmail"`
	ContactPhone    string  `xml:"ContactPhone"`
	Latitude        string  `xml:"Latitude"`
	Longitude       string  `xml:"Longitude"`
	Description     string  `xml:"Description"`
	Category        string  `xml:"Category"`
	OperationType   string  `xml:"OperationType"`
	Price           int64   `xml:"Price"`
	Rooms           string  `xml:"Rooms"`
	Square          float32 `xml:"Square"`
	BalconyOrLoggia string  `xml:"BalconyOrLoggia"`
	KitchenSpace    float32 `xml:"KitchenSpace"`
	ViewFromWindows string  `xml:"ViewFromWindows"`
	CeilingHeight   string  `xml:"CeilingHeight"`
	LivingSpace     float32 `xml:"LivingSpace"`
	Decoration      string  `xml:"Decoration"`
	DealType        string  `xml:"DealType"`
	RoomType        struct {
		Text   string `xml:",chardata"`
		Option string `xml:"Option"`
	} `xml:"RoomType"`
	Status           string `xml:"Status"`
	Floor            int64  `xml:"Floor"`
	Floors           int64  `xml:"Floors"`
	HouseType        string `xml:"HouseType"`
	MarketType       string `xml:"MarketType"`
	PropertyRights   string `xml:"PropertyRights"`
	NewDevelopmentId string `xml:"NewDevelopmentId"`
	Images           struct {
		Image []struct {
			URL string `xml:"url,attr"`
		} `xml:"Image"`
	} `xml:"Images"`
}

func (f *AvitoFeed) Get(url string) (err error) {
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

func (f *AvitoFeed) Check() (errs []error) {

	if len(f.Ad) == 0 {
		errs = append(errs, errors.New("feed is empty"))
		return errs
	}
	if len(f.Ad) <= 10 {
		errs = append(errs, fmt.Errorf("feed contains only %v items", len(f.Ad)))
		return errs
	}

	for idx, lot := range f.Ad {
		if lot.ID == "" {
			errs = append(errs, fmt.Errorf("field ID is empty. Position: %v", idx))
		}
		if lot.Description == "" {
			errs = append(errs, fmt.Errorf("field Description is empty. Position: %v", idx))
		}
		if lot.Category == "" {
			errs = append(errs, fmt.Errorf("field Category is empty. Position: %v", idx))
		}
		if lot.Price == 0 {
			errs = append(errs, fmt.Errorf("field Price is empty. Position: %v", idx))
		}
		if lot.OperationType == "" {
			errs = append(errs, fmt.Errorf("field OperationType is empty. Position: %v", idx))
		}
		if lot.MarketType == "" {
			errs = append(errs, fmt.Errorf("field MarketType is empty. Position: %v", idx))
		}
		if lot.HouseType == "" {
			errs = append(errs, fmt.Errorf("field HouseType is empty. Position: %v", idx))
		}
		if lot.Floor == 0 {
			errs = append(errs, fmt.Errorf("field Floor is empty. Position: %v", idx))
		}
		if lot.Floors == 0 {
			errs = append(errs, fmt.Errorf("field Floors is empty. Position: %v", idx))
		}
		if lot.Rooms == "" {
			errs = append(errs, fmt.Errorf("field Rooms is empty. Position: %v", idx))
		}
		if lot.Square == 0 {
			errs = append(errs, fmt.Errorf("field Square is empty. Position: %v", idx))
		}
		if lot.Status == "" {
			errs = append(errs, fmt.Errorf("field Status is empty. Position: %v", idx))
		}
		if lot.NewDevelopmentId == "" {
			errs = append(errs, fmt.Errorf("field NewDevelopmentId is empty. Position: %v", idx))
		}
		if lot.PropertyRights == "" {
			errs = append(errs, fmt.Errorf("field PropertyRights is empty. Position: %v", idx))
		}
		if lot.Decoration == "" {
			errs = append(errs, fmt.Errorf("field Decoration is empty. Position: %v", idx))
		}

	}
	return errs
}

type AvitoDevelopments struct {
	Region []AvitoRegion `xml:"Region"`
}

type AvitoRegion struct {
	Name string      `xml:"name,attr"`
	City []AvitoCity `xml:"City"`
}

type AvitoCity struct {
	Name   string        `xml:"name,attr"`
	Object []AvitoObject `xml:"Object"`
}

type AvitoObject struct {
	ID        string       `xml:"id,attr"`
	Name      string       `xml:"name,attr"`
	Address   string       `xml:"address,attr"`
	Developer string       `xml:"developer,attr"`
	Housing   []AvitoHouse `xml:"Housing"`
}

type AvitoHouse struct {
	ID      string `xml:"id,attr"`
	Name    string `xml:"name,attr"`
	Address string `xml:"address,attr"`
}

func (f *AvitoFeed) GetDevelopments() (developments AvitoDevelopments, err error) {
	url := "https://autoload.avito.ru/format/New_developments.xml"
	resp, err := GetResponse(url)
	defer resp.Body.Close()
	if err != nil {
		return developments, err
	}

	err = statusCodeHandler(resp)
	if err != nil {
		return developments, err
	}
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return developments, err
	}

	err = xml.Unmarshal(responseBody, &developments)
	if err != nil {
		return developments, err
	}
	//TODO Исправить значение f.LastModified
	return developments, err
}
