package price_placements_feeds

import (
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"time"
)

type DomclickFeed struct {
	LastModified time.Time
	XMLName      xml.Name `xml:"complexes"`
	Complex      struct {
		ID        string `xml:"id"`
		Name      string `xml:"name"`
		Latitude  string `xml:"latitude"`
		Longitude string `xml:"longitude"`
		Address   string `xml:"address"`
		Images    struct {
			Image []string `xml:"image"`
		} `xml:"images"`
		DescriptionMain struct {
			Title string `xml:"title"`
			Text  string `xml:"text"`
		} `xml:"description_main"`
		Infrastructure struct {
			Parking      string `xml:"parking"`
			Security     string `xml:"security"`
			FencedArea   string `xml:"fenced_area"`
			SportsGround string `xml:"sports_ground"`
			Playground   string `xml:"playground"`
			School       string `xml:"school"`
			Kindergarten string `xml:"kindergarten"`
		} `xml:"infrastructure"`
		ProfitsMain struct {
			ProfitMain []struct {
				Title string `xml:"title"`
				Text  string `xml:"text"`
				Image string `xml:"image"`
			} `xml:"profit_main"`
		} `xml:"profits_main"`
		ProfitsSecondary struct {
			ProfitSecondary []struct {
				Title string `xml:"title"`
				Text  string `xml:"text"`
				Image string `xml:"image"`
			} `xml:"profit_secondary"`
		} `xml:"profits_secondary"`
		Buildings struct {
			Building []struct {
				ID            string `xml:"id"`
				Fz214         string `xml:"fz_214"`
				Name          string `xml:"name"`
				Floors        int64  `xml:"floors"`
				BuildingState string `xml:"building_state"`
				BuiltYear     int64  `xml:"built_year"`
				ReadyQuarter  int64  `xml:"ready_quarter"`
				BuildingType  string `xml:"building_type"`
				Image         string `xml:"image"`
				Flats         struct {
					Flat []Flat `xml:"flat"`
				} `xml:"flats"`
			} `xml:"building"`
		} `xml:"buildings"`
		SalesInfo struct {
			SalesPhone              string `xml:"sales_phone"`
			ResponsibleOfficerPhone string `xml:"responsible_officer_phone"`
			SalesAddress            string `xml:"sales_address"`
			SalesLatitude           string `xml:"sales_latitude"`
			SalesLongitude          string `xml:"sales_longitude"`
			Timezone                string `xml:"timezone"`
			WorkDays                struct {
				WorkDay []struct {
					Day     string `xml:"day"`
					OpenAt  string `xml:"open_at"`
					CloseAt string `xml:"close_at"`
				} `xml:"work_day"`
			} `xml:"work_days"`
		} `xml:"sales_info"`
		Developer struct {
			ID    string `xml:"id"`
			Name  string `xml:"name"`
			Phone string `xml:"phone"`
			Site  string `xml:"site"`
			Logo  string `xml:"logo"`
		} `xml:"developer"`
	} `xml:"complex"`
}

type Flat struct {
	FlatID      string  `xml:"flat_id"`
	Apartment   string  `xml:"apartment"`
	Floor       int64   `xml:"floor"`
	Room        *int64  `xml:"room"`
	Plan        string  `xml:"plan"`
	Balcony     string  `xml:"balcony"`
	Renovation  string  `xml:"renovation"`
	Price       int64   `xml:"price"`
	Area        float32 `xml:"area"`
	LivingArea  float32 `xml:"living_area"`
	KitchenArea float32 `xml:"kitchen_area"`
	RoomsArea   struct {
		Area []string `xml:"area"`
	} `xml:"rooms_area"`
	Bathroom     string `xml:"bathroom"`
	HousingType  string `xml:"housing_type"`
	Decoration   int64  `xml:"decoration"`
	ReadyHousing string `xml:"ready_housing"`
}

func (f *DomclickFeed) Get(url string) (err error) {
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

func (f *DomclickFeed) Check() (results []string) {
	if len(f.Complex.Buildings.Building) == 0 {
		results = append(results, emptyFeed)
		return results
	}
	if f.Complex.ID == "" {
		results = append(results, fmt.Sprintf("field Complex.ID is empty"))
	}

	if f.Complex.Name == "" {
		results = append(results, fmt.Sprintf("field Complex.Name is empty"))
	}

	if f.Complex.Address == "" {
		results = append(results, fmt.Sprintf("field Complex.Address is empty"))
	}

	if f.Complex.Latitude == "" {
		results = append(results, fmt.Sprintf("field Complex.Latitude is empty"))
	}

	if f.Complex.Longitude == "" {
		results = append(results, fmt.Sprintf("field Complex.Longitude is empty"))
	}

	for idx, image := range f.Complex.Images.Image {
		if image == "" {
			results = append(results, fmt.Sprintf("field Complex.Images.Image[%d] is empty", idx))
		}
	}

	if f.Complex.DescriptionMain.Title == "" {
		results = append(results, fmt.Sprintf("field Complex.DescriptionMain.Title is empty"))
	}

	if f.Complex.DescriptionMain.Text == "" {
		results = append(results, fmt.Sprintf("field Complex.DescriptionMain.Text is empty"))
	}

	for idx, profit := range f.Complex.ProfitsMain.ProfitMain {
		if profit.Title == "" {
			results = append(results, fmt.Sprintf("field Complex.ProfitsMain.ProfitMain[%d].Title is empty", idx))
		}
		if profit.Text == "" {
			results = append(results, fmt.Sprintf("field Complex.ProfitsMain.ProfitMain[%d].Text is empty", idx))
		}
		if profit.Image == "" {
			results = append(results, fmt.Sprintf("field Complex.ProfitsMain.ProfitMain[%d].Image is empty", idx))
		}
	}

	for pos, building := range f.Complex.Buildings.Building {
		if building.ID == "" {
			results = append(results, fmt.Sprintf("field Complex.Buildings.Building[%d].ID is empty", pos))
		}
		if building.Fz214 == "" {
			results = append(results, fmt.Sprintf("field Complex.Buildings.Building[%s].Fz214 is empty", building.ID))
		}
		if building.Name == "" {
			results = append(results, fmt.Sprintf("field Complex.Buildings.Building[%s].Name is empty", building.ID))
		}
		if building.Floors == 0 {
			results = append(results, fmt.Sprintf("field Complex.Buildings.Building[%s].Floors is empty", building.ID))
		}
		if building.BuildingState == "" {
			results = append(results, fmt.Sprintf("field Complex.Buildings.Building[%s].BuildingState is empty", building.ID))
		}
		if building.BuiltYear == 0 {
			results = append(results, fmt.Sprintf("field Complex.Buildings.Building[%s].BuiltYear is empty", building.ID))
		}
		if building.ReadyQuarter == 0 {
			results = append(results, fmt.Sprintf("field Complex.Buildings.Building[%d].ReadyQuarter is empty", pos))
		}
		if building.BuildingType == "" {
			results = append(results, fmt.Sprintf("field Complex.Buildings.Building[%s].BuildingType is empty", building.ID))
		}

		if building.BuiltYear < int64(time.Now().Year()) && building.BuildingState == "unfinished" {
			results = append(results, fmt.Sprintf("BuildingState == unfinished for %v. InternalID: %v", building.BuiltYear, building.ID))
		}

		for idx, lot := range building.Flats.Flat {
			if lot.FlatID == "" {
				results = append(results, fmt.Sprintf("Field Flats.FlatID is empty. Position: %v", idx))
			}
			if lot.Floor == 0 {
				results = append(results, fmt.Sprintf("Field Flats.Floor is empty. InternalID: %v", lot.FlatID))
			}
			if lot.Room == nil {
				results = append(results, fmt.Sprintf("Field Flats.Room is empty. InternalID: %v", lot.FlatID))
			}
			if lot.Plan == "" {
				results = append(results, fmt.Sprintf("Field Flats.Flat.Plan is empty. InternalID: %v", lot.FlatID))
			}
			if lot.Balcony == "" {
				results = append(results, fmt.Sprintf("Field Flats.Flat.Balcony is empty. InternalID: %v", lot.FlatID))
			}
			if lot.Price == 0 {
				results = append(results, fmt.Sprintf("Field Flats.Flat.Price is empty. InternalID: %v", lot.FlatID))
			}
			if lot.Area == 0 {
				results = append(results, fmt.Sprintf("Field Flats.Flat.Area is empty. InternalID: %v", lot.FlatID))
			}
			if lot.LivingArea == 0 {
				results = append(results, fmt.Sprintf("Field Flats.Flat.LivingArea is empty. InternalID: %v", lot.FlatID))
				for i, room := range lot.RoomsArea.Area {
					if room == "" {
						results = append(results, fmt.Sprintf("Field Flats.Flat.RoomsArea.Area[%v] is empty. InternalID: %v", i, lot.FlatID))
					}
				}
			}
			if lot.KitchenArea == 0 {
				results = append(results, fmt.Sprintf("Field Flats.Flat.KitchenArea is empty. InternalID: %v", lot.FlatID))
			}
			if lot.Bathroom == "" {
				results = append(results, fmt.Sprintf("Field Flats.Flat.Bathroom is empty. InternalID: %v", lot.FlatID))
			}
			if lot.Floor > building.Floors {
				results = append(results, fmt.Sprintf("Field Flats.Flat.Floor is bigger than building.Floors. InternalID: %v", lot.FlatID))
			}
		}
	}

	if f.Complex.SalesInfo.SalesPhone == "" {
		results = append(results, fmt.Sprintf("field Complex.SalesInfo.SalesPhone is empty"))
	}

	if f.Complex.SalesInfo.SalesAddress == "" {
		results = append(results, fmt.Sprintf("field Complex.SalesInfo.SalesAddress is empty"))
	}

	if f.Complex.SalesInfo.SalesLatitude == "" {
		results = append(results, fmt.Sprintf("field Complex.SalesInfo.SalesLatitude is empty"))
	}

	if f.Complex.SalesInfo.SalesLongitude == "" {
		results = append(results, fmt.Sprintf("field Complex.SalesInfo.SalesLongitude is empty"))
	}

	if f.Complex.Developer.Name == "" {
		results = append(results, fmt.Sprintf("field Complex.Developer.Name is empty"))
	}

	if f.Complex.Developer.Phone == "" {
		results = append(results, fmt.Sprintf("field Complex.Developer.Phone is empty"))
	}

	if f.Complex.Developer.Site == "" {
		results = append(results, fmt.Sprintf("field Complex.Developer.Site is empty"))
	}

	if f.Complex.Developer.Logo == "" {
		results = append(results, fmt.Sprintf("field Complex.Developer.Logo is empty"))
	}

	return results
}
