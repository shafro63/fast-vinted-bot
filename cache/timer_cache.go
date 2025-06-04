package cache

import (
	"fast-vinted-bot/utils"
	"log/slog"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/joho/godotenv"
)

//Cache for stop a monitoring link or channel

var _ = godotenv.Load()

func getInterval() time.Duration {
	interval, err := strconv.Atoi(os.Getenv("REFRESH_RATE_TIME"))
	if err != nil {
		slog.Error("unable to convert refresh rate entry", "error", err)
		os.Exit(1)
	}
	return time.Duration(interval) * time.Second
}

var interval = getInterval()

// This will send a tick every xx seconds to refresh the catalog
var generalTicker = time.Tick(interval)

// The timer cache is made of :
// *The channel for the ticks will be sent
// * The duration of the monitor session
// *The The users's sessions (To stop the ticks when an user delete a link from his channel)

var TimerCache = &TimerSessions{
	Sessions: make(map[string]*Timer),
}

type TimerSessions struct {
	Sessions map[string]*Timer
	mu       sync.Mutex
}

type Timer struct {
	TickerChannel chan struct{} // Each monitor will have a channel where the tick will be sent (refresh rate time)
	Duration      time.Duration // Duration of every monitors
}

// Helper functions for the timer cache

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

// The function who will send a tick every xx seconds
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
