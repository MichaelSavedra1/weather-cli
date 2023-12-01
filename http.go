package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

func getRequest(url string, maxRetries int, retryInterval time.Duration) (int, []byte, error) {
	// Retry strategy enabled - Params for this have been defined as constants to avoid assignment
	// on each time this func is called
	for i := 0; i < maxRetries; i++ {

		// Make HTTP request
		res, err := http.Get(url)
		if err != nil {
			continue // Retry on error
		}

		// Read response body
		body, err := ioutil.ReadAll(res.Body)
		res.Body.Close()

		// Check HTTP status code
		if res.StatusCode >= 500 {
			fmt.Printf("Received status code %d, retrying...\n", res.StatusCode)
			time.Sleep(retryInterval)
			continue // Retry on 5xx status codes
		}

		// Process successful
		return res.StatusCode, body, nil
	}

	return 0, nil, fmt.Errorf("Max retries reached")
}
