package apicalls

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"gochopit/utils"
)

func FetchItem(api string, c *utils.AuthCookie) (*utils.Item, error) {

	client := &http.Client{}

	req, err := http.NewRequest("GET", api, nil)
	if err != nil {
		return nil, fmt.Errorf("error while making request :%v", err)
	}

	for k, v := range baseHeaders {
		req.Header.Set(k, v)
	}
	req.Header.Set("Cookie", c.Accesstoken+"; "+c.Refreshtoken)

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error while getting response: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("HTTP %d: %v", resp.StatusCode, string(body))
	}

	defer resp.Body.Close()

	Data := &utils.Item{}

	err = json.NewDecoder(resp.Body).Decode(Data)
	if err != nil {
		return nil, fmt.Errorf("failed to decode JSON: %w", err)
	}
	return Data, nil

}
