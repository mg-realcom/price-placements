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
	Image      []struct {
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
	BuildingType     string      `xml:"building-type"`
	Mortgage         string      `xml:"mortgage"`
	BuildingState    string      `xml:"building-state"`
	Lift             string      `xml:"lift"`
	BathroomUnit     string      `xml:"bathroom-unit"`
	YandexBuildingID int64       `xml:"yandex-building-id"`
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

	if len(f.Offer) == 0 {
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

		if _, ok := tags["plan"]; !ok {
			results = append(results, fmt.Sprintf("tag 'plan' for image is not found. InternalID: %v", lot.InternalID))
		}

		if _, ok := tags["floor-plan"]; !ok {
			results = append(results, fmt.Sprintf("tag 'floor-plan' for image is not found. InternalID: %v", lot.InternalID))
		}

		if lot.Type == "" {
			results = append(results, fmt.Sprintf("field Type is empty. InternalID: %v", lot.InternalID))
		}
		if lot.PropertyType == "" {
			results = append(results, fmt.Sprintf("field PropertyType is empty. InternalID: %v", lot.InternalID))
		}
		if lot.CreationDate == "" {
			results = append(results, fmt.Sprintf("field CreationDate is empty. InternalID: %v", lot.InternalID))
		}
		if lot.Location.Country == "" {
			results = append(results, fmt.Sprintf("field Location.Country is empty. InternalID: %v", lot.InternalID))
		}
		if lot.Location.Address == "" {
			results = append(results, fmt.Sprintf("field Location.Address is empty. InternalID: %v", lot.InternalID))
		}
		if lot.SalesAgent.Phone == "" {
			results = append(results, fmt.Sprintf("field SalesAgent.Phone is empty. InternalID: %v", lot.InternalID))
		}
		if lot.SalesAgent.Category == "" {
			results = append(results, fmt.Sprintf("field SalesAgent.Category is empty. InternalID: %v", lot.InternalID))
		}
		if lot.DealStatus == "" {
			results = append(results, fmt.Sprintf("field DealStatus is empty. InternalID: %v", lot.InternalID))
		}
		if lot.Price.Value == 0 {
			results = append(results, fmt.Sprintf("field Price.Value is empty. InternalID: %v", lot.InternalID))
		}
		if lot.Price.Currency == "" {
			results = append(results, fmt.Sprintf("field Price.Currency is empty. InternalID: %v", lot.InternalID))
		}
		if lot.Area.Value == 0 {
			results = append(results, fmt.Sprintf("field Area.Value is empty. InternalID: %v", lot.InternalID))
		}
		if lot.Area.Unit == "" {
			results = append(results, fmt.Sprintf("field Area.Unit is empty. InternalID: %v", lot.InternalID))
		}
		if lot.LivingSpace.Value == 0 && lot.OpenPlan != "1" {
			results = append(results, fmt.Sprintf("field LivingSpace.Value is empty. InternalID: %v", lot.InternalID))
		}
		if lot.Rooms == 0 {
			results = append(results, fmt.Sprintf("field Rooms is empty. InternalID: %v", lot.InternalID))
		}
		if lot.NewFlat == "" {
			results = append(results, fmt.Sprintf("field NewFlat is empty. InternalID: %v", lot.InternalID))
		}
		if lot.Floor == 0 {
			results = append(results, fmt.Sprintf("field Floor is empty. InternalID: %v", lot.InternalID))
		}
		if lot.FloorsTotal == 0 {
			results = append(results, fmt.Sprintf("field FloorsTotal is empty. InternalID: %v", lot.InternalID))
		}
		if lot.BuildingName == "" {
			results = append(results, fmt.Sprintf("field BuildingName is empty. InternalID: %v", lot.InternalID))
		}
		if lot.YandexBuildingID == 0 {
			results = append(results, fmt.Sprintf("field YandexBuildingID is empty. InternalID: %v", lot.InternalID))
		}
		if lot.BuildingState == "" {
			results = append(results, fmt.Sprintf("field BuildingState is empty. InternalID: %v", lot.InternalID))
		}
		if lot.BuiltYear == 0 {
			results = append(results, fmt.Sprintf("field BuiltYear is empty. InternalID: %v", lot.InternalID))
		}
		if lot.ReadyQuarter == 0 {
			results = append(results, fmt.Sprintf("field ReadyQuarter is empty. InternalID: %v", lot.InternalID))
		}

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
