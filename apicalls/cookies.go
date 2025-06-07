package apicalls

import (
	"log/slog"
	"net/http"
	"net/http/cookiejar"
	"sync"
	"time"

	"fast-vinted-bot/utils"
)

var mu sync.Mutex

func GetCookie(rb *utils.RequestBuilder) []*http.Cookie {
	vintedUrl := "https://" + rb.URL.Host
	jar, _ := cookiejar.New(nil)

	client := &http.Client{
		Transport: &http.Transport{
			Proxy:             http.ProxyURL(rb.Proxy),
			DisableKeepAlives: true,
		},
		Timeout: 3 * time.Second,
	}
	client.Jar = jar

	req, err := http.NewRequest(rb.Method, vintedUrl, nil)
	if err != nil {
		slog.Error("unable to create request", "error", err)
		return nil
	}

	req.Header = map[string][]string{
		"User-Agent":      {utils.Ua.GetRandom()},
		"Accept":          {"text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7"},
		"Accept-encoding": {"gzip, deflate, br, zstd"},
		"Accept-language": {"fr"},
	}

	resp, err := client.Do(req)
	if err != nil {
		slog.Error("unable to send request", "error", err)
		return nil
	}
	err = utils.HandleHttpError(client, req, resp)
	if err != nil {
		slog.Error("unable to get cookies", "error", err)
		return nil
	}
	defer resp.Body.Close()

	cookies := client.Jar.Cookies(req.URL)
	if len(cookies) == 0 {
		return nil
	}

	// Empty the jar so the future requests won't include useless cookies
	client.Jar = nil
	slog.Debug("Fetched cookies", "cookies", cookies)
	return cookies
}

// Function to get api access tokens
func FormatedAuthCookie(c []*http.Cookie) *utils.AuthCookie {

	var accesstoken, refreshtoken string
	for _, c := range c {
		if c.Name == "access_token_web" {
			accesstoken = c.String()
		} else if c.Name == "refresh_token_web" {
			refreshtoken = c.String()
		}
	}

	authcookie := &utils.AuthCookie{}

	authcookie.Accesstoken.String = accesstoken
	authcookie.Refreshtoken.String = refreshtoken

	// The accesstoken cookie expires every 2 hours,
	authcookie.Accesstoken.Expires = 110 * time.Minute

	slog.Debug("Cookies formated")
	return authcookie

}

// So it needs to be refreshed with his refreshtoken
func RefreshCookie(rb *utils.RequestBuilder) {
	go func() {
		for {
			accesstoken_life := rb.Cookie.Accesstoken.Expires
			time.Sleep(accesstoken_life)

			mu.Lock()

			rb.Cookie.Accesstoken = rb.Cookie.Refreshtoken

			refreshed_cookie := GetCookie(rb)
			new_cookie := FormatedAuthCookie(refreshed_cookie)

			rb.Cookie = new_cookie

			mu.Unlock()

			slog.Debug("Cookies refreshed")
		}
	}()
}
