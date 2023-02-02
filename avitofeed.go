package price_placements_feeds

import (
	"encoding/xml"
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

func (f *AvitoFeed) Check() (results []string) {

	if len(f.Ad) == 0 {
		results = append(results, emptyFeed)
		return results
	}

	if len(f.Ad) <= 10 {
		results = append(results, fmt.Sprintf("feed contains only %v items", len(f.Ad)))
		return results
	}

	for idx, lot := range f.Ad {
		if lot.ID == "" {
			results = append(results, fmt.Sprintf("field ID is empty. Position: %v", idx))
		}
		if lot.Description == "" {
			results = append(results, fmt.Sprintf("field Description is empty. Position: %v", idx))
		}
		if lot.Category == "" {
			results = append(results, fmt.Sprintf("field Category is empty. Position: %v", idx))
		}
		if lot.Price == 0 {
			results = append(results, fmt.Sprintf("field Price is empty. Position: %v", idx))
		}
		if lot.OperationType == "" {
			results = append(results, fmt.Sprintf("field OperationType is empty. Position: %v", idx))
		}
		if lot.MarketType == "" {
			results = append(results, fmt.Sprintf("field MarketType is empty. Position: %v", idx))
		}
		if lot.HouseType == "" {
			results = append(results, fmt.Sprintf("field HouseType is empty. Position: %v", idx))
		}
		if lot.Floor == 0 {
			results = append(results, fmt.Sprintf("field Floor is empty. Position: %v", idx))
		}
		if lot.Floors == 0 {
			results = append(results, fmt.Sprintf("field Floors is empty. Position: %v", idx))
		}
		if lot.Rooms == "" {
			results = append(results, fmt.Sprintf("field Rooms is empty. Position: %v", idx))
		}
		if lot.Square == 0 {
			results = append(results, fmt.Sprintf("field Square is empty. Position: %v", idx))
		}
		if lot.Status == "" {
			results = append(results, fmt.Sprintf("field Status is empty. Position: %v", idx))
		}
		if lot.NewDevelopmentId == "" {
			results = append(results, fmt.Sprintf("field NewDevelopmentId is empty. Position: %v", idx))
		}
		if lot.PropertyRights == "" {
			results = append(results, fmt.Sprintf("field PropertyRights is empty. Position: %v", idx))
		}
		if lot.Decoration == "" {
			results = append(results, fmt.Sprintf("field Decoration is empty. Position: %v", idx))
		}

	}
	return results
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
