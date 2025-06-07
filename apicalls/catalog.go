package apicalls

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"fast-vinted-bot/utils"
)

func FetchCatalogItems(rb *utils.RequestBuilder) ([]utils.CatalogItem, error) {
	client := &http.Client{
		Transport: &http.Transport{
			Proxy:             http.ProxyURL(rb.Proxy),
			DisableKeepAlives: true,
		},
		Timeout: 3 * time.Second,
	}

	api := fmt.Sprintf("%s://%s%s/catalog/items?%s", rb.URL.Scheme, rb.URL.Host, rb.URL.Path, rb.URL.RawQuery)

	req, _ := http.NewRequest(rb.Method, api, nil)
	req.Header = map[string][]string{
		"Cache-Control":   {"no-cache"},
		"Pragma":          {"no-cache"},
		"User-Agent":      {utils.Ua.GetRandom()},
		"Accept":          {"text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7"},
		"Accept-encoding": {"gzip, deflate, br, zstd"},
		"Accept-language": {"fr"},
		"Cookie":          {rb.Cookie.Accesstoken.String},
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("unable to send request : %v ", err)
	}
	err = utils.HandleHttpError(client, req, resp)
	if err != nil {
		return nil, fmt.Errorf("fetch catalog item : %v ", err)
	}

	defer resp.Body.Close()

	var Data utils.ItemsResp
	var CatalogItems []utils.CatalogItem

	err = json.NewDecoder(resp.Body).Decode(&Data)
	if err != nil {
		return nil, fmt.Errorf("failed to decode JSON: %w", err)
	}

	CatalogItems = append(CatalogItems, Data.Items...)

	return CatalogItems, nil
}
