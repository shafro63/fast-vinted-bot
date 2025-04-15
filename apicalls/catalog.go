package apicalls

import (
	"encoding/json"
	"fmt"
	"net/http"

	"gochopit/utils"
)

func FetchCatalogItems(api string, c *utils.AuthCookie) ([]utils.CatalogItems, error) {
	client := &http.Client{}

	req, _ := http.NewRequest("GET", api, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	req.Header.Set("X-Requested-With", "XMLHttpRequest") // Simule une requête AJAX
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Cookie", c.Accesstoken+"; "+c.Refreshtoken)

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Erreur lors de la requête :%v", err)
		return nil, err
	}
	defer resp.Body.Close()

	var Data utils.ItemsResp
	var CatalogItems []utils.CatalogItems

	err = json.NewDecoder(resp.Body).Decode(&Data)
	if err != nil {
		return nil, fmt.Errorf("failed to decode JSON: %w", err)
	}

	CatalogItems = append(CatalogItems, Data.Items...)

	return CatalogItems, nil
}
