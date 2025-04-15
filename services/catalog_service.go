package services

import (
	"fmt"
	"net/url"
	"time"

	"gochopit/apicalls"
	"gochopit/utils"
)

var Catalogchan = make(chan []utils.CatalogItems, 1)

var ItemsChan = make(chan []*utils.Item)

var lastId int64

var Duration = 60 * time.Second

var concurrency = 0

func SetFetchInterval(function, interval time.Time) {

}

func SetCatalogApi(u *url.URL) string {
	return fmt.Sprintf("%s://%s%s/catalog/items?%s", u.Scheme, u.Host, u.Path, u.RawQuery)
}

func SetItemsApi(u *url.URL, id int64) string {
	return fmt.Sprintf("%s://%s%s/items/%d", u.Scheme, u.Host, u.Path, id)
}

func NewItems(items []utils.CatalogItems) []utils.CatalogItems {
	var newItems []utils.CatalogItems
	var i0 = items[0]
	if lastId != 0 {
		for _, item := range items {
			if item.ID == lastId {
				break
			}
			newItems = append(newItems, item)
		}
		lastId = i0.ID
	} else {
		lastId = i0.ID
	}
	return newItems
}

func FetchCatalogAtInterval(c *utils.AuthCookie, u *url.URL, interval time.Duration) error {
	catalogApi := SetCatalogApi(u)
	ticker := time.NewTicker(interval)

	done := make(chan bool)
	go func() {
		time.Sleep(Duration)
		done <- true
	}()

	for {
		select {
		case t := <-ticker.C:
			f, err := apicalls.FetchCatalogItems(catalogApi, c)
			if err != nil {
				return fmt.Errorf("error while fetching catalog at interval %v : %v", t, err)
			}
			Catalogchan <- f
		case <-done:
			close(Catalogchan)
			return nil
		}
	}
}

func Fetchclean() {
	for items := range Catalogchan {
		newitems := NewItems(items)
		for i := len(newitems) - 1; i >= 0; i-- {
			item := newitems[i]
			fmt.Printf("Id : %d | Title : %s | Visible : %v\n", item.ID, item.Title, item.IsVisible)
		}
	}

}

func FetchAndSendItems(u *url.URL, c *utils.AuthCookie, ids []int64, t *time.Ticker) error {
	if len(ids) == 0 {
		return nil
	}
	//ticker := time.NewTicker(1 * time.Second)

	ch := make(chan struct{}, concurrency)

	for _, itemid := range ids {
		itemApi := SetItemsApi(u, itemid)

		go func(api string) {
			ch <- struct{}{}
			defer func() { <-ch }()

			item, err := apicalls.FetchItem(api, c)
			if err != nil {
				fmt.Printf("Failed to fetch Item : %v", err)
			}

			fmt.Printf(`Id : %d
Titre : %s
Marque : %s`, item.ID, item.Title, item.Brand)
		}(itemApi)
	}
	return nil
}
