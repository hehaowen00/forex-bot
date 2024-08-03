package main

import (
	"bufio"
	"encoding/json"
	"log"
	"net/http"
	oanda "oanda-api/api"
)

func StartStream(params *oanda.GetCandlesRes) error {
	url := "https://stream-fxpractice.oanda.com/v3/accounts/101-011-11236738-001/pricing/stream?instruments=EUR_USD&account_id=101-011-11236738-001"

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}

	// req.Header.Set("Authorization", "Bearer "+config.AuthToken)
	resp, err := http.DefaultClient.Do(req)

	reader := bufio.NewReader(resp.Body)

	for {
		line, _, err := reader.ReadLine()
		if err != nil {
			log.Println(err)
			continue
		}

		data := json.RawMessage{}

		err = json.Unmarshal(line, &data)
		if err != nil {
			return err
		}

		log.Println(string(data))
	}

	return nil
}
