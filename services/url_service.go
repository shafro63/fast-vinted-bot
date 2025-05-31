package services

import (
	"fmt"
	"net/url"
	"strings"
)

func ParsedUrl(link string) (*url.URL, error) {
	parsedUrl, err := url.Parse(link)
	if err != nil {
		return nil, fmt.Errorf("invalid URL: %v", err)
	}
	if !strings.Contains(parsedUrl.Host, "vinted.") {
		return nil, fmt.Errorf("invalid URL: this is not a Vinted URL")
	}
	if !strings.HasPrefix(parsedUrl.Host, "www.") {
		parsedUrl.Host = "www." + parsedUrl.Host
	}
	if !strings.HasPrefix(parsedUrl.Path, "/catalog") {
		return nil, fmt.Errorf("invalid URL: You should only put urls with the /catalog route")
	}

	params := parsedUrl.Query()
	if len(params) == 0 {
		return nil, fmt.Errorf("no filters in url. You have to select filters")
	}

	newParams := url.Values{}

	for key, values := range params {
		newKey := key
		if strings.HasSuffix(key, "[]") {
			newKey = strings.ReplaceAll(key, "catalog", "catalog_ids")
			newKey = strings.ReplaceAll(newKey, "[]", "")
		}
		newParams[newKey] = values
	}
	newParams.Set("order", "newest_first")

	return &url.URL{
		Scheme:   "https",
		Host:     parsedUrl.Host,
		Path:     "/web/api/core",
		RawQuery: newParams.Encode(),
	}, nil
}
