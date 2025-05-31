package utils

import (
	"fmt"
	"net/http"
)

func CheckStatusCode(resp *http.Response) error {

	switch resp.StatusCode {
	case 200:
		return nil
	case 301, 302:
		return fmt.Errorf("redirection : %v", resp.Status)
	case 400:
		return fmt.Errorf("bad Request : %v", resp.Status)
	case 401:
		return fmt.Errorf("unauthorized : %v", resp.Status)
	case 403:
		return fmt.Errorf("access denied : %v", resp.Status)
	case 404:
		return fmt.Errorf("page not found : %v", resp.Status)
	case 500:
		return fmt.Errorf("server error : %v", resp.Status)
	default:
		return fmt.Errorf("request error :%v", resp.Status)
	}
}
