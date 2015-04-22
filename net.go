package github

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type RateLimitError struct {
	What string
	When time.Time
}

func (e RateLimitError) Error() string {
	return fmt.Sprintf("%v: %v", e.When, e.What)
}

func httpGetBody(url string) ([]byte, error) {

	resp, err := http.Get(url)
	defer resp.Body.Close()
	if err != nil {
		return []byte{}, err
	}

	if resp.StatusCode == 403 { // Forbidden
		return []byte{}, RateLimitError{"the rate limit was exceeded", time.Now()}
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return body, err
	}

	return body, nil
}
