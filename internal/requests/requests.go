package requests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"
)

var (
	client = http.Client{
		Timeout: 90 * time.Second,
	}
)

func GetData(url string, v interface{}) error {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(v)
	if err != nil {
		return err
	}

	return nil
}

func PostData(url string, v interface{}) (*http.Response, error) {
	responseBytes, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}

	body := bytes.NewBuffer(responseBytes)

	postReq, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		return nil, err
	}
	postReq.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(postReq)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
