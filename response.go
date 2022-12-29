package price_placements_feeds

import (
	"fmt"
	"net/http"
)

func GetResponse(url string) (response *http.Response, err error) {
	response, err = http.Get(url)
	if response.StatusCode != 200 {
		return response, fmt.Errorf("feed not availible. Status:%s", response.Status)
	}
	if err != nil {
		return response, fmt.Errorf("can't get feed. Error:%w", err)
	}
	return response, err
}
