package price_placements_feeds

import (
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"time"
)

type RealtyFeed struct {
	LastModified   time.Time
	Xmlns          string  `xml:"xmlns,attr"`
	GenerationDate string  `xml:"generation-date"`
	Offer          []Offer `xml:"offer"`
}

type Offer struct {
	InternalID string `xml:"internal-id,attr"`

	Image []struct {
		Tag string `xml:"tag,attr"`
	} `xml:"image"`

	Type           string   `xml:"type"`
	PropertyType   string   `xml:"property-type"`
	Category       string   `xml:"category"`
	URL            string   `xml:"url"`
	WindowView     string   `xml:"window-view"`
	CeilingHeight  []string `xml:"ceiling-height"`
	Description    string   `xml:"description"`
	CreationDate   string   `xml:"creation-date"`
	Vas            []vas    `xml:"vas"`
	LastUpdateDate string   `xml:"last-update-date"`
	ExpireDate     string   `xml:"expire-date"`
	Location       struct {
		Country      string `xml:"country"`
		Region       string `xml:"region"`
		Address      string `xml:"address"`
		LocalityName string `xml:"locality-name"`
		Latitude     string `xml:"latitude"`
		Longitude    string `xml:"longitude"`
		Direction    string `xml:"direction"`
		Distance     string `xml:"distance"`
		Metro        struct {
			Name            string `xml:"name"`
			TimeOnTransport string `xml:"time-on-transport"`
			TimeOnFoot      string `xml:"time-on-foot"`
		} `xml:"metro"`
	} `xml:"location"`
	SalesAgent struct {
		Category     string `xml:"category"`
		Organization string `xml:"organization"`
		Phone        string `xml:"phone"`
	} `xml:"sales-agent"`
	Price struct {
		Value    float32 `xml:"value"`
		Currency string  `xml:"currency"`
	} `xml:"price"`
	NewFlat          string      `xml:"new-flat"`
	DealStatus       string      `xml:"deal-status"`
	BuiltYear        int64       `xml:"built-year"`
	ReadyQuarter     int64       `xml:"ready-quarter"`
	Area             Value       `xml:"area"`
	RoomSpace        []Value     `xml:"room-space"`
	LivingSpace      Value       `xml:"living-space"`
	KitchenSpace     Value       `xml:"kitchen-space"`
	Renovation       string      `xml:"renovation"`
	Rooms            int64       `xml:"rooms"`
	RubbishChute     string      `xml:"rubbish-chute"`
	FloorsTotal      int64       `xml:"floors-total"`
	Floor            int64       `xml:"floor"`
	BuildingName     string      `xml:"building-name"`
	VillageName      string      `xml:"village-name"`
	BuildingType     string      `xml:"building-type"`
	Mortgage         string      `xml:"mortgage"`
	BuildingState    string      `xml:"building-state"`
	Lift             string      `xml:"lift"`
	BathroomUnit     string      `xml:"bathroom-unit"`
	YandexBuildingID int64       `xml:"yandex-building-id"`
	YandexVillageID  int64       `xml:"yandex-village-id"`
	YandexHouseID    CustomInt64 `xml:"yandex-house-id"`
	BuildingSection  string      `xml:"building-section"`
	Balcony          string      `xml:"balcony"`
	OpenPlan         string      `xml:"open-plan"`
}

type Value struct {
	Value float32 `xml:"value"`
	Unit  string  `xml:"unit"`
}

type vas struct {
	Text      string `xml:",chardata"`
	StartTime string `xml:"start-time,attr"`
	Schedule  string `xml:"schedule,attr"`
}

func (f *RealtyFeed) Get(url string) (err error) {
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
	if time.Time.IsZero(f.LastModified) {
		f.LastModified, err = time.Parse(time.RFC3339Nano, f.GenerationDate)
		if err != nil {
			return err
		}
	}
	return nil
}

func (f *RealtyFeed) Check() (results []string) {

	if len(f.Offer) < 2 {
		results = append(results, emptyFeed)
		return results
	}

	for idx, lot := range f.Offer {
		if lot.InternalID == "" {
			results = append(results, fmt.Sprintf("field InternalID is empty. Position: %v", idx))
		}
		tags := make(map[string]bool)
		for _, image := range lot.Image {
			if tags[image.Tag] {
				continue
			}
			tags[image.Tag] = true
		}

		if lot.Type == "" {
			results = append(results, fmt.Sprintf("tag 'Type'  is not found. InternalID: %v", lot.InternalID))
		}

		id := lot.InternalID

		if lot.VillageName != "" {
			checkStringWithID(id, "offer", "VillageName", lot.VillageName, &results)
			checkZeroWithID(id, "offer", "YandexVillageID", int(lot.YandexVillageID), &results)
		} else {
			checkStringWithID(id, "offer", "BuildingName", lot.BuildingName, &results)
			checkZeroWithID(id, "offer", "YandexBuildingID", int(lot.YandexBuildingID), &results)
			checkStringWithID(id, "offer", "NewFlat", lot.NewFlat, &results)
			if lot.LivingSpace.Value == 0 && lot.OpenPlan != "1" {
				results = append(results, fmt.Sprintf("field LivingSpace.Value is empty. InternalID: %v", lot.InternalID))
			}
			if _, ok := tags["floor-plan"]; !ok {
				results = append(results, fmt.Sprintf("tag 'floor-plan' for image is not found. InternalID: %v", lot.InternalID))
			}
			if _, ok := tags["plan"]; !ok {
				results = append(results, fmt.Sprintf("tag 'plan' for image is not found. InternalID: %v", lot.InternalID))
			}
		}

		checkStringWithID(id, "offer", "Type", lot.Type, &results)
		checkStringWithID(id, "offer", "PropertyType", lot.PropertyType, &results)
		checkStringWithID(id, "offer", "CreationDate", lot.CreationDate, &results)
		checkStringWithID(id, "offer.Location", "Country", lot.Location.Country, &results)
		checkStringWithID(id, "offer.Location", "Address", lot.Location.Address, &results)
		checkStringWithID(id, "offer.SalesAgent", "Phone", lot.SalesAgent.Phone, &results)
		checkStringWithID(id, "offer.SalesAgent", "Category", lot.SalesAgent.Category, &results)
		checkStringWithID(id, "offer", "DealStatus", lot.DealStatus, &results)
		checkZeroWithID(id, "offer.Price", "Value", lot.Price.Value, &results)
		checkStringWithID(id, "offer.Price", "Currency", lot.Price.Currency, &results)
		checkZeroWithID(id, "offer.Area", "Value", lot.Area.Value, &results)
		checkStringWithID(id, "offer.Area", "Unit", lot.Area.Unit, &results)
		checkZeroWithID(id, "offer", "Rooms", int(lot.Rooms), &results)

		checkZeroWithID(id, "offer", "Floor", int(lot.Floor), &results)
		checkZeroWithID(id, "offer", "FloorsTotal", int(lot.FloorsTotal), &results)
		checkStringWithID(id, "offer", "BuildingState", lot.BuildingState, &results)
		checkZeroWithID(id, "offer", "BuiltYear", int(lot.BuiltYear), &results)
		checkZeroWithID(id, "offer", "ReadyQuarter", int(lot.ReadyQuarter), &results)

		if lot.BuiltYear < int64(time.Now().Year()) && lot.BuildingState == "unfinished" {
			results = append(results, fmt.Sprintf("BuildingState == unfinished for %v. InternalID: %v", lot.BuiltYear, lot.InternalID))
		}
		if lot.Floor > lot.FloorsTotal {
			results = append(results, fmt.Sprintf("field Floor is bigger than FloorsTotal. InternalID: %v", lot.InternalID))
		}
		if int64(len(lot.RoomSpace)) > lot.Rooms {
			results = append(results, fmt.Sprintf("field RoomSpace contains more values than Rooms. InternalID: %v", lot.InternalID))
		}
		if len(lot.Image) < 3 {
			results = append(results, fmt.Sprintf("field Image contains '%v' items. InternalID: %v", len(lot.Image), lot.InternalID))
		}
	}
	return results
}
