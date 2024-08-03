package oandaapi

import (
	"bytes"
	"encoding/json"
	"net/http"
)

func newJSONRequest[Req any](method string, url string, req *Req) (*http.Request, error) {
	data, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	return http.NewRequest(method, url, bytes.NewReader(data))
}

func postJSON[Req any](client *http.Client, url string, payload *Req) (*http.Response, error) {
	data, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	resp, err := client.Post(url, "application/json", bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func unmarshalJSON[T any](resp *http.Response) (*T, error) {
	var res T

	dec := json.NewDecoder(resp.Body)

	err := dec.Decode(&res)
	if err != nil {
		return nil, err
	}

	return &res, nil
}
