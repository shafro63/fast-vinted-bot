package apicalls

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"fast-vinted-bot/utils"
)

func FetchItem(rb *utils.RequestBuilder, id int) (*utils.Item, error) {
	client := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(rb.Proxy),
		},
	}

	api := fmt.Sprintf("%s://%s%s/items/%d", rb.URL.Scheme, rb.URL.Host, rb.URL.Path, id)
	req, err := http.NewRequest(rb.Method, api, nil)
	if err != nil {
		return nil, fmt.Errorf("error while making request :%v", err)
	}

	req.Header.Set("Cookie", rb.Cookie.Accesstoken.String)

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
