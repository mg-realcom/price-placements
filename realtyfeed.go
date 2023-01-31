package price_placements_feeds

import (
	"encoding/xml"
	"errors"
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
		Text string `xml:",chardata"`
		Tag  string `xml:"tag,attr"`
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
	NewFlat          string  `xml:"new-flat"`
	DealStatus       string  `xml:"deal-status"`
	BuiltYear        int64   `xml:"built-year"`
	ReadyQuarter     int64   `xml:"ready-quarter"`
	Area             Value   `xml:"area"`
	RoomSpace        []Value `xml:"room-space"`
	LivingSpace      Value   `xml:"living-space"`
	KitchenSpace     Value   `xml:"kitchen-space"`
	Renovation       string  `xml:"renovation"`
	Rooms            int64   `xml:"rooms"`
	RubbishChute     string  `xml:"rubbish-chute"`
	FloorsTotal      int64   `xml:"floors-total"`
	Floor            int64   `xml:"floor"`
	BuildingName     string  `xml:"building-name"`
	BuildingType     string  `xml:"building-type"`
	Mortgage         string  `xml:"mortgage"`
	BuildingState    string  `xml:"building-state"`
	Lift             string  `xml:"lift"`
	BathroomUnit     string  `xml:"bathroom-unit"`
	YandexBuildingID int64   `xml:"yandex-building-id"`
	YandexHouseID    int64   `xml:"yandex-house-id"`
	BuildingSection  string  `xml:"building-section"`
	Balcony          string  `xml:"balcony"`
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
	if time.Time.IsZero(f.LastModified) {
		f.LastModified, err = time.Parse(time.RFC3339Nano, f.GenerationDate)
		if err != nil {
			return err
		}
	}
	return nil
}

func (f *RealtyFeed) Check() (errs []error) {

	if len(f.Offer) == 0 {
		errs = append(errs, errors.New("feed is empty"))
		return errs
	}
	for idx, lot := range f.Offer {
		if lot.InternalID == "" {
			errs = append(errs, fmt.Errorf("field InternalID is empty. Position: %v", idx))
		}
		if lot.Type == "" {
			errs = append(errs, fmt.Errorf("field Type is empty. InternalID: %v", lot.InternalID))
		}
		if lot.PropertyType == "" {
			errs = append(errs, fmt.Errorf("field PropertyType is empty. InternalID: %v", lot.InternalID))
		}
		if lot.CreationDate == "" {
			errs = append(errs, fmt.Errorf("field CreationDate is empty. InternalID: %v", lot.InternalID))
		}
		if lot.Location.Country == "" {
			errs = append(errs, fmt.Errorf("field Location.Country is empty. InternalID: %v", lot.InternalID))
		}
		if lot.SalesAgent.Phone == "" {
			errs = append(errs, fmt.Errorf("field SalesAgent.Phone is empty. InternalID: %v", lot.InternalID))
		}
		if lot.SalesAgent.Category == "" {
			errs = append(errs, fmt.Errorf("field SalesAgent.Category is empty. InternalID: %v", lot.InternalID))
		}
		if lot.DealStatus == "" {
			errs = append(errs, fmt.Errorf("field DealStatus is empty. InternalID: %v", lot.InternalID))
		}
		if lot.Price.Value == 0 {
			errs = append(errs, fmt.Errorf("field Price.Value is empty. InternalID: %v", lot.InternalID))
		}
		if lot.Price.Currency == "" {
			errs = append(errs, fmt.Errorf("field Price.Currency is empty. InternalID: %v", lot.InternalID))
		}
		if lot.Area.Value == 0 {
			errs = append(errs, fmt.Errorf("field Area.Value is empty. InternalID: %v", lot.InternalID))
		}
		if lot.Area.Unit == "" {
			errs = append(errs, fmt.Errorf("field Area.Unit is empty. InternalID: %v", lot.InternalID))
		}
		if lot.NewFlat == "" {
			errs = append(errs, fmt.Errorf("field NewFlat is empty. InternalID: %v", lot.InternalID))
		}
		if lot.Floor == 0 {
			errs = append(errs, fmt.Errorf("field Floor is empty. InternalID: %v", lot.InternalID))
		}
		if lot.FloorsTotal == 0 {
			errs = append(errs, fmt.Errorf("field FloorsTotal is empty. InternalID: %v", lot.InternalID))
		}
		if lot.BuildingName == "" {
			errs = append(errs, fmt.Errorf("field BuildingName is empty. InternalID: %v", lot.InternalID))
		}
		if lot.YandexBuildingID == 0 {
			errs = append(errs, fmt.Errorf("field YandexBuildingID is empty. InternalID: %v", lot.InternalID))
		}
		if lot.YandexHouseID == 0 {
			errs = append(errs, fmt.Errorf("field YandexHouseID is empty. InternalID: %v", lot.InternalID))
		}
		if lot.BuildingState == "" {
			errs = append(errs, fmt.Errorf("field BuildingState is empty. InternalID: %v", lot.InternalID))
		}
		if lot.BuiltYear == 0 {
			errs = append(errs, fmt.Errorf("field BuiltYear is empty. InternalID: %v", lot.InternalID))
		}
		if lot.ReadyQuarter == 0 {
			errs = append(errs, fmt.Errorf("field ReadyQuarter is empty. InternalID: %v", lot.InternalID))
		}
	}
	return errs
}
