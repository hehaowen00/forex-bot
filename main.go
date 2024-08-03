package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	oanda "oanda-api/api"
	"os"
	"strconv"
	"time"

	go_env "github.com/hehaowen00/go-env"

	_ "github.com/mattn/go-sqlite3"
)

const (
	PRACTICE_URL = "https://api-fxpractice.oanda.com"
)

type Config struct {
	AccountID string `env:"ACCOUNT_ID"`
	AuthToken string `env:"AUTH_TOKEN"`
}

func main() {
	config, err := go_env.LoadFile[Config](".env")
	if err != nil {
		log.Fatalln(err)
	}

	api := oanda.NewAPIClient(
		http.DefaultClient,
		config.AccountID,
		config.AuthToken,
		PRACTICE_URL,
	)

	majors := []string{
		"EUR_USD", "USD_JPY", "GBP_USD", "USD_CHF", "USD_CAD", "AUD_USD", "NZD_USD",
	}

	if false {
		res, err := api.GetInstruments(&oanda.GetInstrumentsReq{
			Instruments: majors,
		})
		if err != nil {
			log.Fatalln(err)
		}

		data, err := json.MarshalIndent(res, "", "  ")
		if err != nil {
			log.Fatalln(err)
		}

		log.Println(os.WriteFile("instruments.json", data, 0777))
	}

	if true {
		for _, major := range majors {
			log.Println("major", major)
			instrument := major
			granularity := oanda.GRANULARITY_M1

			since, _ := time.Parse(time.RFC3339, "2024-08-02T20:59:00Z")

			candlesRes, err := api.GetCandles(&oanda.GetCandlesReq{
				Instrument:  instrument,
				Granularity: granularity,
				Price:       "MBA",
				Count:       5000,
				From:        strconv.Itoa(int(since.Unix())),
			})
			if err != nil {
				log.Fatalln(err)
			}

			data, err := json.MarshalIndent(candlesRes, "", "  ")
			if err != nil {
				log.Fatalln(err)
			}

			log.Println(os.WriteFile(instrument+"_"+granularity+"_1.json", data, 0777))

			db, err := sql.Open("sqlite3", "")
			if err != nil {
				log.Panic(err)
			}

			db.Exec(``)
		}
	}
}
