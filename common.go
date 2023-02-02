package price_placements_feeds

import (
	"fmt"
	"net/http"
)

const (
	emptyFeed string = "feed is empty"
)

func statusCodeHandler(resp *http.Response) error {
	if resp == nil {
		return fmt.Errorf("не могу получить ответ сервера")
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("не могу получить ответ сервера. Statuscode: %v", resp.StatusCode)
	}
	return nil
}
