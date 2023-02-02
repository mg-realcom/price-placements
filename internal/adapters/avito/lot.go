package avito

import (
	"fmt"
	placementsFeeds "github.com/zfullio/price-placements"
	"price-placements-service/internal/domain/entity"
	"regexp"
	"strconv"
	"strings"
)

type lotRepository struct {
	client placementsFeeds.AvitoFeed
}

func NewLotRepository() *lotRepository {
	return &lotRepository{client: placementsFeeds.AvitoFeed{}}
}

func (lr lotRepository) Get(url string) (lots []entity.Lot, err error) {
	err = lr.client.Get(url)
	if err != nil {
		return lots, fmt.Errorf("не могу разобрать фид: %w", err)
	}
	developments, err := lr.getObjectMap()
	if err != nil {
		return lots, fmt.Errorf("не могу получить словарь DevelopmentID: %w", err)

	}
	for i := 0; i < len(lr.client.Ad); i++ {
		onlyDigitInt, err := phoneNumberToInt(lr.client.Ad[i].ContactPhone)
		if err != nil {
			return nil, fmt.Errorf("не могу сконвертировать номер в число")
		}
		name, ok := developments[lr.client.Ad[i].NewDevelopmentId]
		if !ok {
			return lots, fmt.Errorf("не могу сопоставить DevelopmentID: %s", lr.client.Ad[i].NewDevelopmentId)
		}
		lot := entity.Lot{
			ID:     lr.client.Ad[i].ID,
			Object: optimizeObject(name),
			Phone:  onlyDigitInt,
		}
		lots = append(lots, lot)
	}
	return lots, err
}

func (lr lotRepository) getObjectMap() (map[string]string, error) {
	objectMap := make(map[string]string, 0)
	developments, err := lr.client.GetDevelopments()
	if err != nil {
		return objectMap, err
	}
	for _, region := range developments.Region {
		for _, city := range region.City {
			for _, object := range city.Object {
				if _, ok := objectMap[object.ID]; ok != true {
					objectMap[object.ID] = object.Name
				}
			}
		}
	}
	return objectMap, nil
}

func optimizeObject(str string) string {
	result := strings.ToLower(str)
	result = strings.ReplaceAll(result, "\"", "")
	result = strings.ReplaceAll(result, "«", "")
	result = strings.ReplaceAll(result, "»", "")
	result = strings.TrimPrefix(result, "жк ")
	result = strings.TrimPrefix(result, "ск ")
	result = strings.ReplaceAll(result, "апарт-комплекс", "")
	result = strings.ReplaceAll(result, "сити-комплекс", "")
	result = strings.TrimSpace(result)
	return result
}

func phoneNumberToInt(str string) (phone int, err error) {
	re := regexp.MustCompile(`\D+`)
	res := re.ReplaceAllString(str, "")
	re = regexp.MustCompile(`([0-9]{11})`)
	onlyDigitStr := string(re.Find([]byte(res)))
	phone, err = strconv.Atoi(onlyDigitStr)
	if err != nil {
		return phone, fmt.Errorf("не могу сконвертировать номер в число")
	}
	return phone, err
}

func (lr lotRepository) Validate(url string) (results []string, err error) {
	err = lr.client.Get(url)
	if err != nil {
		return results, fmt.Errorf("не могу разобрать фид: %w", err)
	}
	results = append(results, lr.client.Check()...)
	return results, err
}
