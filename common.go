package price_placements_feeds

import (
	"fmt"
	"net/http"
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
