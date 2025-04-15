package apicalls

import (
	"log/slog"
	"net/http"
	"net/http/cookiejar"
	"net/url"

	"gochopit/utils"
)

var baseHeaders = map[string]string{
	"User-Agent":      "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36",
	"Accept":          "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7",
	"Accept-encoding": "gzip, deflate, br, zstd",
	"Accept-language": "fr",
}

func GetCookie(link string) []*http.Cookie {
	parsedUrl, _ := url.Parse(link)
	vintedUrl := "https://" + parsedUrl.Host
	jar, _ := cookiejar.New(nil)
	client := &http.Client{Jar: jar}

	req, err := http.NewRequest("GET", vintedUrl, nil)
	if err != nil {
		slog.Debug("Erreur lors de la création de la requête", "erreur", err)
		return nil
	}

	for k, v := range baseHeaders {
		req.Header.Set(k, v)
	}
	req.Header.Set("Referer", parsedUrl.Host)

	resp, err := client.Do(req)
	if err != nil {
		slog.Debug("Erreur lors de la réception de la réponse", "erreur", err)
		return nil
	}
	defer resp.Body.Close()

	cookies := client.Jar.Cookies(req.URL)
	if len(cookies) == 0 {
		slog.Debug("aucun cookie trouvé", "erreur", err)
		return nil
	}
	return cookies
}

func FormatedAuthCookie(c []*http.Cookie) *utils.AuthCookie {
	var accesstoken, refreshtoken string
	for _, c := range c {
		if c.Name == "access_token_web" {
			accesstoken = c.String()
		} else if c.Name == "refresh_token_web" {
			refreshtoken = c.String()
		}
	}

	return &utils.AuthCookie{
		Accesstoken:  accesstoken,
		Refreshtoken: refreshtoken,
	}

}

func RefreshCookie(link string) *utils.AuthCookie {
	cookie := GetCookie(link)
	return FormatedAuthCookie(cookie)
}
