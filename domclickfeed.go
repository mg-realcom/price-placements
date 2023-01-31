package price_placements_feeds

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"log"
	"time"
)

type DomclickFeed struct {
	LastModified time.Time
	XMLName      xml.Name `xml:"complexes"`
	Text         string   `xml:",chardata"`
	Complex      struct {
		Text      string `xml:",chardata"`
		ID        string `xml:"id"`
		Name      string `xml:"name"`
		Latitude  string `xml:"latitude"`
		Longitude string `xml:"longitude"`
		Address   string `xml:"address"`
		Images    struct {
			Text  string   `xml:",chardata"`
			Image []string `xml:"image"`
		} `xml:"images"`
		DescriptionMain struct {
			Chardata string `xml:",chardata"`
			Title    string `xml:"title"`
			Text     string `xml:"text"`
		} `xml:"description_main"`
		Infrastructure struct {
			Text         string `xml:",chardata"`
			Parking      string `xml:"parking"`
			Security     string `xml:"security"`
			FencedArea   string `xml:"fenced_area"`
			SportsGround string `xml:"sports_ground"`
			Playground   string `xml:"playground"`
			School       string `xml:"school"`
			Kindergarten string `xml:"kindergarten"`
		} `xml:"infrastructure"`
		ProfitsMain struct {
			Text       string `xml:",chardata"`
			ProfitMain []struct {
				Chardata string `xml:",chardata"`
				Title    string `xml:"title"`
				Text     string `xml:"text"`
				Image    string `xml:"image"`
			} `xml:"profit_main"`
		} `xml:"profits_main"`
		ProfitsSecondary struct {
			Text            string `xml:",chardata"`
			ProfitSecondary []struct {
				Chardata string `xml:",chardata"`
				Title    string `xml:"title"`
				Text     string `xml:"text"`
				Image    string `xml:"image"`
			} `xml:"profit_secondary"`
		} `xml:"profits_secondary"`
		Buildings struct {
			Text     string `xml:",chardata"`
			Building []struct {
				Text          string `xml:",chardata"`
				ID            string `xml:"id"`
				Fz214         string `xml:"fz_214"`
				Name          string `xml:"name"`
				Floors        string `xml:"floors"`
				BuildingState string `xml:"building_state"`
				BuiltYear     string `xml:"built_year"`
				ReadyQuarter  string `xml:"ready_quarter"`
				BuildingType  string `xml:"building_type"`
				Image         string `xml:"image"`
				Flats         struct {
					Text string `xml:",chardata"`
					Flat []Flat `xml:"flat"`
				} `xml:"flats"`
			} `xml:"building"`
		} `xml:"buildings"`
		SalesInfo struct {
			Text                    string `xml:",chardata"`
			SalesPhone              string `xml:"sales_phone"`
			ResponsibleOfficerPhone string `xml:"responsible_officer_phone"`
			SalesAddress            string `xml:"sales_address"`
			SalesLatitude           string `xml:"sales_latitude"`
			SalesLongitude          string `xml:"sales_longitude"`
			Timezone                string `xml:"timezone"`
			WorkDays                struct {
				Text    string `xml:",chardata"`
				WorkDay []struct {
					Text    string `xml:",chardata"`
					Day     string `xml:"day"`
					OpenAt  string `xml:"open_at"`
					CloseAt string `xml:"close_at"`
				} `xml:"work_day"`
			} `xml:"work_days"`
		} `xml:"sales_info"`
		Developer struct {
			Text  string `xml:",chardata"`
			ID    string `xml:"id"`
			Name  string `xml:"name"`
			Phone string `xml:"phone"`
			Site  string `xml:"site"`
			Logo  string `xml:"logo"`
		} `xml:"developer"`
	} `xml:"complex"`
}

type Flat struct {
	Text        string  `xml:",chardata"`
	FlatID      string  `xml:"flat_id"`
	Apartment   string  `xml:"apartment"`
	Floor       int64   `xml:"floor"`
	Room        string  `xml:"room"`
	Plan        string  `xml:"plan"`
	Balcony     string  `xml:"balcony"`
	Renovation  string  `xml:"renovation"`
	Price       int64   `xml:"price"`
	Area        float32 `xml:"area"`
	LivingArea  float32 `xml:"living_area"`
	KitchenArea float32 `xml:"kitchen_area"`
	RoomsArea   struct {
		Text string   `xml:",chardata"`
		Area []string `xml:"area"`
	} `xml:"rooms_area"`
	Bathroom     string `xml:"bathroom"`
	HousingType  string `xml:"housing_type"`
	Decoration   int64  `xml:"decoration"`
	ReadyHousing string `xml:"ready_housing"`
}

func (f *DomclickFeed) Get(url string) (err error) {
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

func (f *DomclickFeed) Check() (errs []error) {
	if len(f.Complex.Buildings.Building) == 0 {
		errs = append(errs, errors.New("feed is empty"))
		return errs
	}
	if f.Complex.ID == "" {
		errs = append(errs, fmt.Errorf("field Complex.ID is empty"))
	}

	if f.Complex.Name == "" {
		errs = append(errs, fmt.Errorf("field Complex.Name is empty"))
	}

	if f.Complex.Address == "" {
		errs = append(errs, fmt.Errorf("field Complex.Address is empty"))
	}

	if f.Complex.DescriptionMain.Title == "" {
		errs = append(errs, fmt.Errorf("field Complex.DescriptionMain.Title is empty"))
	}

	if f.Complex.DescriptionMain.Text == "" {
		errs = append(errs, fmt.Errorf("field Complex.DescriptionMain.Text is empty"))
	}

	for _, building := range f.Complex.Buildings.Building {
		for idx, lot := range building.Flats.Flat {
			if lot.FlatID == "" {
				errs = append(errs, fmt.Errorf("BuildingId: %s. Field Flats.FlatID is empty. Position: %v \n", building.ID, idx))
			}
			if lot.FlatID == "" {
				errs = append(errs, fmt.Errorf("BuildingId: %s. Field Flats.Flat.FlatID is empty. Position: %v \n", building.ID, idx))
			}
			if lot.Room == "" {
				errs = append(errs, fmt.Errorf("BuildingId: %s. Field Flats.Flat.Room is empty. Position: %v \n", building.ID, idx))
			}
			if lot.Price == 0 {
				errs = append(errs, fmt.Errorf("BuildingId: %s. Field Flats.Flat.Price is empty. Position: %v \n", building.ID, idx))
			}
			if lot.Area == 0 {
				errs = append(errs, fmt.Errorf("BuildingId: %s. Field Flats.Flat.Area is empty. Position: %v \n", building.ID, idx))
			}
			if lot.Area == 0 {
				errs = append(errs, fmt.Errorf("BuildingId: %s. Field Flats.Flat.Area is empty. Position: %v \n", building.ID, idx))
			}
		}
	}

	return errs
}
