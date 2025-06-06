package utils

import (
	"log/slog"
	"net/http"
	"net/url"
	"os"

	"github.com/joho/godotenv"
	randua "github.com/lib4u/fake-useragent"
)

var _ = godotenv.Load()

func getProxy() *url.URL {
	link := os.Getenv("PROXY_URL")
	proxy, err := url.Parse(link)
	if err != nil {
		slog.Error("unable to get proxy url", "error", err)
		os.Exit(1)
	}
	return proxy
}

func getUserAgent() *randua.UserAgent {
	ua, err := randua.New()
	if err != nil {
		slog.Error("unable to load user agents", "error", err)
		os.Exit(1)
	}
	return ua
}

var Ua = getUserAgent()

var Headers = map[string][]string{
	"User-Agent":      {Ua.GetRandom()}, //random user agent
	"Accept":          {"text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7"},
	"Accept-encoding": {"gzip, deflate, br, zstd"},
	"Accept-language": {"fr"},
}

// Request builder for sending requests
func (rb *RequestBuilder) New() {
	rb.Proxy = getProxy()
	rb.Client = &http.Client{
		Jar: nil,
		Transport: &http.Transport{
			Proxy: http.ProxyURL(rb.Proxy),
		},
	}

}
