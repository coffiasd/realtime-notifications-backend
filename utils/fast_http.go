package utils

import (
	"io/ioutil"
	"net/http"
	"time"
)

func Request(url string,timeout time.Duration) ([]byte,error){
	client := &http.Client{Timeout: timeout}
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return []byte(""),err
	}
	res, err := client.Do(req)
	if err != nil {
		return []byte(""),err
	}

	//close handle
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return []byte(""),err
	}

	return body,nil
}