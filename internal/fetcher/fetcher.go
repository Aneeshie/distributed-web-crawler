package fetcher

import (
	"io"
	"net/http"
	"time"
)

var client = &http.Client{
	Timeout: 10 * time.Second,
}

func Fetch(url string) (int, []byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return 0, nil, err
	}

	req.Header.Set("User-Agent","WebCrawler/1.0" )

	resp,err := client.Do(req)
	if err != nil {
		return 0, nil, err
	}

	defer resp.Body.Close()
	//for now ill be using readAll for production based systems preferrably ill limit this 
	body, err := io.ReadAll(resp.Body)
	if err != nil  {
		return resp.StatusCode, nil, err
	}

	return resp.StatusCode, body, err

}

