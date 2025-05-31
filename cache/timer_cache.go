package cache

import (
	"log/slog"
	"sync"
	"time"

	"fast-vinted-bot/utils"
)

var interval = 2 * time.Second

var generalTicker = time.Tick(interval)

var TimerCache = &TimerSessions{
	Sessions: make(map[string]*Timer),
}

type TimerSessions struct {
	Sessions map[string]*Timer
	mu       sync.Mutex
}

type Timer struct {
	TickerChannel chan struct{}
	Duration      time.Duration
}

func (t *TimerSessions) SetTimer(linkName string, timer *Timer) {
	TimerCache.mu.Lock()
	defer TimerCache.mu.Unlock()

	TimerCache.Sessions[linkName] = timer
}

func (t *TimerSessions) DeleteTimer(linkName string) {
	TimerCache.mu.Lock()
	defer TimerCache.mu.Unlock()

	delete(TimerCache.Sessions, linkName)
}

func (t *TimerSessions) DeleteAllTimersInChannel(data *utils.DiscordUserData) {
	TimerCache.mu.Lock()
	defer TimerCache.mu.Unlock()

	session := DataCache.GetMonitoringChannel(data)
	if session == nil {
		return
	}
	for linkName := range session.Links {
		delete(TimerCache.Sessions, linkName)
	}

}

func LaunchTicker() {
	go func() {
		for range generalTicker {
			TimerCache.mu.Lock()
			for _, t := range TimerCache.Sessions {
				t.TickerChannel <- struct{}{}
			}
			TimerCache.mu.Unlock()
		}
	}()
	slog.Info("Ticker launched")
}
