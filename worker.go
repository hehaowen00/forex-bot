package main

import (
	"log"
	oandaapi "oanda-api/api"
	"time"
)

type Worker struct {
	api *oandaapi.APIClient
}

func NewWorker(api *oandaapi.APIClient) *Worker {
	return &Worker{
		api: api,
	}
}

func (w *Worker) Run() {
	err := w.api.StartStream(&oandaapi.GetCandlesReq{}, func(sr *oandaapi.StreamRes) error {
		return nil
	})
	if err != nil {
		log.Println(err)
	}

	for {
		time.Sleep(time.Second * 5)
	}
}
