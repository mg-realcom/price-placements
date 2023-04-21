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
	Price       float32 `xml:"price"`
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
	if len(f.Complex.Buildings.Building) < 2 {
		results = append(results, emptyFeed)
		return results
	}
	path := "Complex"
	checkString(path, "ID", f.Complex.ID, &results)
	checkString(path, "Name", f.Complex.Name, &results)
	checkString(path, "Address", f.Complex.Address, &results)
	checkString(path, "Latitude", f.Complex.Latitude, &results)
	checkString(path, "Longitude", f.Complex.Longitude, &results)

	for idx, image := range f.Complex.Images.Image {
		checkStringWithPos(idx, "Complex.Images.Image", "Image", image, &results)
	}

	path = "Complex.DescriptionMain"
	checkString(path, "Title", f.Complex.DescriptionMain.Title, &results)
	checkString(path, "Text", f.Complex.DescriptionMain.Text, &results)

	for idx, profit := range f.Complex.ProfitsMain.ProfitMain {
		path := "Complex.ProfitsMain.ProfitMain"
		checkStringWithPos(idx, path, "Title", profit.Title, &results)
		checkStringWithPos(idx, path, "Text", profit.Text, &results)
		checkStringWithPos(idx, path, "Image", profit.Image, &results)
	}

	for pos, building := range f.Complex.Buildings.Building {
		path := "Complex.Buildings.Building"
		checkStringWithPos(pos, path, "ID", building.ID, &results)
		checkStringWithID(building.ID, path, "Fz214", building.Fz214, &results)
		checkStringWithID(building.ID, path, "Name", building.Name, &results)
		checkZeroWithID(building.ID, path, "Floors", int(building.Floors), &results)
		checkStringWithID(building.ID, path, "BuildingState", building.BuildingState, &results)
		checkZeroWithID(building.ID, path, "BuiltYear", int(building.BuiltYear), &results)
		checkZeroWithID(building.ID, path, "ReadyQuarter", int(building.ReadyQuarter), &results)
		checkStringWithID(building.ID, path, "BuildingType", building.BuildingType, &results)

		if building.BuiltYear < int64(time.Now().Year()) && building.BuildingState == "unfinished" {
			results = append(results, fmt.Sprintf("BuildingState == unfinished for %v. InternalID: %v", building.BuiltYear, building.ID))
		}

		f.checkLots(building.Flats.Flat, int(building.Floors), &results)
	}

	path = "Complex.SalesInfo"
	checkString(path, "SalesPhone", f.Complex.SalesInfo.SalesPhone, &results)
	checkString(path, "SalesAddress", f.Complex.SalesInfo.SalesAddress, &results)
	checkString(path, "SalesLatitude", f.Complex.SalesInfo.SalesLatitude, &results)
	checkString(path, "SalesLongitude", f.Complex.SalesInfo.SalesLongitude, &results)

	path = "Complex.Developer"
	checkString(path, "Name", f.Complex.Developer.Name, &results)
	checkString(path, "Phone", f.Complex.Developer.Phone, &results)
	checkString(path, "Site", f.Complex.Developer.Site, &results)
	checkString(path, "Logo", f.Complex.Developer.Logo, &results)

	return results
}

func (f *DomclickFeed) checkLots(flats []Flat, floors int, results *[]string) {
	for idx, lot := range flats {
		path := "Flats.Flat"
		checkStringWithPos(idx, path, "FlatID", lot.FlatID, results)
		checkZeroWithID(lot.FlatID, path, "Floor", int(lot.Floor), results)
		if lot.Room == nil {
			*results = append(*results, fmt.Sprintf("Field Flats.Room is empty. InternalID: %v", lot.FlatID))
		}
		checkStringWithID(lot.FlatID, path, "Plan", lot.Plan, results)
		checkStringWithID(lot.FlatID, path, "Balcony", lot.Balcony, results)
		checkZeroWithID(lot.FlatID, path, "Price", lot.Price, results)
		checkZeroWithID(lot.FlatID, path, "Area", lot.Area, results)
		isOk := checkZeroWithID(lot.FlatID, path, "LivingArea", lot.LivingArea, results)
		if !isOk {
			for i, room := range lot.RoomsArea.Area {
				if room == "" {
					*results = append(*results, fmt.Sprintf("Field Flats.Flat.RoomsArea.Area[%v] is empty. InternalID: %v", i, lot.FlatID))
				}
			}
		}

		checkZeroWithID(lot.FlatID, path, "KitchenArea", lot.KitchenArea, results)
		checkStringWithID(lot.FlatID, path, "Bathroom", lot.Bathroom, results)

		if lot.Floor > int64(floors) {
			*results = append(*results, fmt.Sprintf("Field Flats.Flat.Floor is bigger than building.Floors. InternalID: %v", lot.FlatID))
		}
	}
}
