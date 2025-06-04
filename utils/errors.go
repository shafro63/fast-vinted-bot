package utils

import (
	"fmt"
	"log/slog"
	"net/http"
)

// Handle the http errors for debugging
func HandleHttpError(client *http.Client, req *http.Request, resp *http.Response) error {
	max_retries := 3

	switch resp.StatusCode {

	case 200:
		return nil

	case 301, 302:
		return fmt.Errorf("redirection : %v", resp.Status)

	case 400:
		return fmt.Errorf("bad request : %v", resp.Status)

	case 401:
		return fmt.Errorf("unauthorized : %v", resp.Status)

	case 403:
		resp, err := client.Do(req)
		if err != nil {
			return fmt.Errorf("response error : %v", resp.Status)
		}

		retries_count := 0

		for range max_retries {
			resp, err = client.Do(req)
			if err != nil {
				return fmt.Errorf("response error : %v", resp.Status)
			}

			if resp.StatusCode == 200 {
				return nil
			}

			retries_count++
			slog.Debug("403 Error", "retries", retries_count)
		}

		return fmt.Errorf("too many 403 errors : %v", resp.Status)

	case 404:
		return fmt.Errorf("page not found : %v", resp.Status)

	case 500:
		return fmt.Errorf("server error : %v", resp.Status)

	default:
		return fmt.Errorf("request error :%v", resp.Status)
	}
}
