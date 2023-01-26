package avito

import (
	"checkphones/internal/domain/entity"
	"fmt"
	placementsFeeds "github.com/zfullio/price-placements"
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
	for i := 0; i < len(lr.client.Ad); i++ {
		onlyDigitInt, err := phoneNumberToInt(lr.client.Ad[i].ContactPhone)
		if err != nil {
			return nil, fmt.Errorf("не могу сконвертировать номер в число")
		}

		lot := entity.Lot{
			ID:     lr.client.Ad[i].ID,
			Object: optimizeObject(lr.client.Ad[i].Description),
			Phone:  onlyDigitInt,
		}
		lots = append(lots, lot)
	}
	return lots, err
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
	re := regexp.MustCompile(`([0-9]{11})`)
	onlyDigitStr := string(re.Find([]byte(str)))
	phone, err = strconv.Atoi(onlyDigitStr)
	if err != nil {
		return phone, fmt.Errorf("не могу сконвертировать номер в число")
	}
	return phone, err
}