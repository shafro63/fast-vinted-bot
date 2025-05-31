package services

import (
	"context"
	"log/slog"
	"time"

	"fast-vinted-bot/apicalls"
	"fast-vinted-bot/cache"
	"fast-vinted-bot/utils"
)

func NewItems(items []utils.CatalogItem, lastId int64) ([]utils.CatalogItem, int64) {
	var newItems []utils.CatalogItem
	if len(items) == 0 {
		return nil, 0
	}
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
	if len(newItems) == 0 {
		slog.Debug("New items Success", "items", newItems)
	}
	slog.Debug("New items Success")
	return newItems, lastId
}

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
					slog.Debug("error while fetching catalog", "error", err)
					continue
				}
				dataChan <- f
			case <-ctx.Done():
				end := time.Since(start)
				slog.Debug("Fetching Items End", "duration", end)
				return
			case <-stopChan:
				end := time.Since(start)
				slog.Debug("Fetch stopped", "duration", end)
				return
			}
		}
	}()
}
