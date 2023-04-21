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

func (f *AvitoFeed) Check() (results []string) {

	if len(f.Ad) < 2 {
		results = append(results, emptyFeed)
		return results
	}

	if len(f.Ad) <= 10 {
		results = append(results, fmt.Sprintf("feed contains only %v items", len(f.Ad)))
		return results
	}

	for idx, lot := range f.Ad {
		checkStringWithPos(idx, "Ad", "ID", lot.ID, &results)
		id := lot.ID
		checkStringWithID(id, "Ad", "ContactPhone", lot.ContactPhone, &results)
		checkStringWithID(id, "Ad", "Description", lot.Description, &results)
		checkStringWithID(id, "Ad", "Category", lot.Category, &results)
		checkZeroWithID(id, "Ad", "Price", int(lot.Price), &results)
		checkStringWithID(id, "Ad", "OperationType", lot.OperationType, &results)
		checkStringWithID(id, "Ad", "MarketType", lot.MarketType, &results)
		checkStringWithID(id, "Ad", "HouseType", lot.HouseType, &results)
		checkZeroWithID(id, "Ad", "Floor", int(lot.Floor), &results)
		checkZeroWithID(id, "Ad", "Floors", int(lot.Floors), &results)
		checkStringWithID(id, "Ad", "Rooms", lot.Rooms, &results)
		checkZeroWithID(id, "Ad", "Square", lot.Square, &results)

		if lot.LivingSpace == 0 && lot.Rooms != "Студия" {
			results = append(results, fmt.Sprintf("field LivingSpace is empty. InternalID: %v", lot.ID))
		}

		checkStringWithID(id, "Ad", "Status", lot.Status, &results)
		checkStringWithID(id, "Ad", "NewDevelopmentId", lot.NewDevelopmentId, &results)
		checkStringWithID(id, "Ad", "PropertyRights", lot.PropertyRights, &results)
		checkStringWithID(id, "Ad", "Decoration", lot.Decoration, &results)

		if lot.Floor > lot.Floors {
			results = append(results, fmt.Sprintf("field Floor is bigger than Floors. InternalID: %v", lot.ID))
		}
		for idx, image := range lot.Images.Image {
			checkStringWithPos(idx, "Images.Image", "URL", image.URL, &results)
		}

		if len(lot.Images.Image) < 3 || len(lot.Images.Image) > 40 {
			results = append(results, fmt.Sprintf("field Images.Image contains '%v' items. InternalID: %v", len(lot.Images.Image), lot.ID))
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
