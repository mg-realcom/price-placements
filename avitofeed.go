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
	FormatVersion string   `xml:"formatVersion,attr"`
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

func (f *AvitoFeed) Check() (err error) {
	if len(f.Ad) == 0 {
		return errors.New("feed is empty")
	}
	for idx, lot := range f.Ad {
		if lot.ID == "" {
			return fmt.Errorf("field ID is empty. Position: %v", idx)
		}
		if lot.Description == "" {
			return fmt.Errorf("field Description is empty. Position: %v", idx)
		}
		if lot.Category == "" {
			return fmt.Errorf("field Category is empty. Position: %v", idx)
		}
		if lot.Price == 0 {
			return fmt.Errorf("field Price is empty. Position: %v", idx)
		}
		if lot.OperationType == "" {
			return fmt.Errorf("field OperationType is empty. Position: %v", idx)
		}
		if lot.MarketType == "" {
			return fmt.Errorf("field MarketType is empty. Position: %v", idx)
		}
		if lot.HouseType == "" {
			return fmt.Errorf("field HouseType is empty. Position: %v", idx)
		}
		if lot.Floor == 0 {
			return fmt.Errorf("field Floor is empty. Position: %v", idx)
		}
		if lot.Floors == 0 {
			return fmt.Errorf("field Floors is empty. Position: %v", idx)
		}
		if lot.Rooms == "" {
			return fmt.Errorf("field Rooms is empty. Position: %v", idx)
		}
		if lot.Square == 0 {
			return fmt.Errorf("field Square is empty. Position: %v", idx)
		}
		if lot.Status == "" {
			return fmt.Errorf("field Status is empty. Position: %v", idx)
		}
		if lot.NewDevelopmentId == "" {
			return fmt.Errorf("field NewDevelopmentId is empty. Position: %v", idx)
		}
		if lot.PropertyRights == "" {
			return fmt.Errorf("field PropertyRights is empty. Position: %v", idx)
		}
		if lot.Decoration == "" {
			return fmt.Errorf("field Decoration is empty. Position: %v", idx)
		}

	}
	return err
}