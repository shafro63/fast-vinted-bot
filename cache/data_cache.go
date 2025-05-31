package cache

import (
	"fast-vinted-bot/utils"
	"sync"
)

var DataCache = &MiniCache{
	MonitoringChannels: make(map[string]*Session),
	UsersData:          make(map[string]*utils.DiscordUserData),
}

type MiniCache struct {
	MonitoringChannels map[string]*Session               // ChannelID, User's session object
	UsersData          map[string]*utils.DiscordUserData // User choices throught interaction inputs
	Mu                 sync.Mutex
}

type Session struct {
	Links map[string]chan bool // Link's name, Go channels where the ads are sent
}

func (c *MiniCache) GetMonitoringChannel(data *utils.DiscordUserData) *Session {
	c.Mu.Lock()
	defer c.Mu.Unlock()

	session := c.MonitoringChannels[data.ChannelID]
	return session
}

func (c *MiniCache) DeleteMonitoringChannel(data *utils.DiscordUserData) {
	c.Mu.Lock()
	defer c.Mu.Unlock()

	for k, v := range c.MonitoringChannels {
		if k == data.ChannelID {
			for _, ch := range v.Links {
				ch <- false
			}

		}
	}
	delete(c.MonitoringChannels, data.ChannelID)
}

func (c *MiniCache) SetMonitorSession(data *utils.DiscordUserData, ch chan bool) {
	c.Mu.Lock()
	defer c.Mu.Unlock()

	session := c.MonitoringChannels[data.ChannelID]
	if session == nil {
		session = &Session{
			Links: make(map[string]chan bool),
		}
		c.MonitoringChannels[data.ChannelID] = session
	}

	session.Links[data.LinkName] = ch
}

func (c *MiniCache) GetMonitorSession(data *utils.DiscordUserData) chan bool {
	c.Mu.Lock()
	defer c.Mu.Unlock()

	session := c.MonitoringChannels[data.ChannelID]
	if session == nil {
		return nil
	}

	ch := session.Links[data.LinkName]
	if ch == nil {
		return nil
	}
	return ch
}

func (c *MiniCache) DeleteMonitorSession(data *utils.DiscordUserData) {
	c.Mu.Lock()
	defer c.Mu.Unlock()

	session := c.MonitoringChannels[data.ChannelID]
	stopChan := session.Links[data.LinkName]

	close(stopChan)
	delete(session.Links, data.LinkName)
}

func (c *MiniCache) SetUserData(userID string, data *utils.DiscordUserData) {
	c.Mu.Lock()
	defer c.Mu.Unlock()

	c.UsersData[userID] = data
}

func (c *MiniCache) GetUserData(userID string) *utils.DiscordUserData {
	c.Mu.Lock()
	defer c.Mu.Unlock()

	userdata := c.UsersData[userID]
	if userdata == nil {
		return nil
	}
	data := *userdata

	delete(c.UsersData, userID)
	return &data
}
