package oandaapi

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
)

var (
	CURRENCIES = []string{"AUD", "CAD", "CHF", "EUR", "GBP", "JPY", "USD"}
)

type APIClient struct {
	client    *http.Client
	accountID string
	authToken string
	url       string
}

func NewAPIClient(client *http.Client, accountID string, authToken string, url string) *APIClient {
	return &APIClient{
		client:    client,
		accountID: accountID,
		authToken: authToken,
		url:       url,
	}
}

func (api *APIClient) addAuth(req *http.Request) {
	// req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+api.authToken)
}

func (api *APIClient) GetInstruments(
	req *GetInstrumentsReq,
) (*GetInstrumentsRes, error) {
	url := api.url + "/v3/accounts/" + api.accountID + "/instruments"

	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	if len(req.Instruments) > 0 {
		q := request.URL.Query()
		q.Add("instruments", strings.Join(req.Instruments, ","))
		request.URL.RawQuery = q.Encode()
	}

	api.addAuth(request)

	resp, err := api.client.Do(request)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("get instruments error - status code %d %s", resp.StatusCode, string(body))
	}

	res, err := unmarshalJSON[GetInstrumentsRes](resp)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (api *APIClient) GetCandles(req *GetCandlesReq) (*GetCandlesRes, error) {
	url := api.url + "/v3/instruments/" + req.Instrument + "/candles"

	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	q := request.URL.Query()
	q.Set("granularity", req.Granularity)
	q.Set("count", strconv.Itoa(int(req.Count)))
	q.Set("price", req.Price)
	if req.From != "" {
		q.Set("from", req.From)
	}

	request.URL.RawQuery = q.Encode()

	log.Println(request.URL.String())

	api.addAuth(request)

	resp, err := api.client.Do(request)
	if err != nil {
		return nil, err
	}

	res, err := unmarshalJSON[GetCandlesRes](resp)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (api *APIClient) GetOrderBook(req *GetOrderBookReq) (*GetOrderBookRes, error) {
	url := api.url + "/v3/instruments/" + req.Instrument + "/orderBook"

	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := api.client.Do(request)
	if err != nil {
		return nil, err
	}

	res, err := unmarshalJSON[GetOrderBookRes](resp)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (api *APIClient) GetOrders(req *GetOrdersReq) (*GetOrdersRes, error) {
	url := api.url + "/v3/accounts/" + api.accountID + "/orders"

	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	if req != nil {
		q := request.URL.Query()
		if req.Instrument != "" {
			q.Set("instrument", req.Instrument)
		}
		request.URL.RawQuery = q.Encode()
	}

	resp, err := api.client.Do(request)
	if err != nil {
		return nil, err
	}

	res, err := unmarshalJSON[GetOrdersRes](resp)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (api *APIClient) StartStream(params *GetCandlesReq, callback func(*StreamRes) error) error {
	url := "https://stream-fxpractice.oanda.com/v3/accounts/" + api.accountID + "/pricing/stream"

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}

	q := req.URL.Query()
	q.Set("instruments", params.Instrument)

	req.URL.RawQuery = q.Encode()

	req.Header.Set("Authorization", "Bearer "+api.authToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	reader := bufio.NewReader(resp.Body)

	for {
		line, _, err := reader.ReadLine()
		if err != nil {
			log.Println(err)
			continue
		}

		res := StreamRes{}

		data := json.RawMessage{}

		err = json.Unmarshal(line, &data)
		if err != nil {
			return err
		}

		log.Println(string(data))

		err = callback(&res)
		if err != nil {
			return err
		}
	}
}
