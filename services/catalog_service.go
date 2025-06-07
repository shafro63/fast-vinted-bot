package services

import (
	"context"
	"log/slog"
	"time"

	"fast-vinted-bot/apicalls"
	"fast-vinted-bot/cache"
	"fast-vinted-bot/utils"
)

// Fetch the catalog and send the items throught a channel for filtering
func FetchCatalogAtInterval(rb *utils.RequestBuilder, timer *cache.Timer, dataChan chan<- []utils.CatalogItem, stopChan chan bool) {

	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), timer.Duration)
		defer cancel()
		defer close(dataChan)

		start := time.Now()

		for {
			select {
			case <-timer.TickerChannel:
				f, err := apicalls.FetchCatalogItems(rb)
				if err != nil {
					slog.Error("Error FetchCatalogInterval", "error", err, "url", rb.URL)
					continue
				}
				dataChan <- f
			case <-ctx.Done():
				end := time.Since(start)
				slog.Info("Fetching Items End", "duration", end)
				return
			case <-stopChan:
				end := time.Since(start)
				slog.Info("Fetch stopped", "duration", end)
				return
			}
		}
	}()
}

// Get only latest items by filtering by last item id captured
func LatestItems(items []utils.CatalogItem, lastId *int64) []utils.CatalogItem {
	if len(items) == 0 {
		return nil
	}
	newest := items[0].ID
	previous := *lastId
	if previous == 0 {
		*lastId = newest
		return nil
	}

	var result []utils.CatalogItem
	for _, item := range items {
		if item.ID == previous {
			break
		}
		result = append(result, item)
	}
	if len(result) == 0 {
		slog.Debug("no new items")
		return result
	} else {
		slog.Debug("new items !")
	}
	return result
}
