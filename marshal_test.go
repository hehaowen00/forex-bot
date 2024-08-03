package main

import (
	"encoding/json"
	"log"
	oanda "oanda-api/api"
	"os"
	"testing"
	"time"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
)

func TestCandlestickMarshal(t *testing.T) {
	contents, err := os.ReadFile("USD_JPY_S30_1.json")
	if err != nil {
		t.FailNow()
	}

	res := oanda.GetCandlesRes{}

	err = json.Unmarshal(contents, &res)
	if err != nil {
		t.FailNow()
	}

	loc, err := time.LoadLocation("Australia/Brisbane")
	if err != nil {
		t.FailNow()
	}

	ts := res.Candles[0].Time
	log.Println(ts, ts.UnixMilli(), ts.In(loc))

	ts = res.Candles[len(res.Candles)-1].Time
	log.Println(ts, ts.UnixMilli(), ts.In(loc))

	var avgDiff time.Duration
	count := 0

	for i, candle := range res.Candles {
		if i == 0 {
			continue
		}

		ts0 := res.Candles[i-1].Time
		ts1 := candle.Time

		diff := ts1.Sub(ts0)

		if diff > time.Second*30 {
			log.Println("diff (", i-1, i, diff, ")", ts0.In(loc), ts1.In(loc))
		} else {
			avgDiff += diff
			count += 1
		}
	}

	log.Println(len(res.Candles), count, (avgDiff)/time.Duration(count), time.Second*30*time.Duration(len(res.Candles)))
}

func TestCandlestickPattern(t *testing.T) {
	contents, err := os.ReadFile("USD_JPY_S30_1.json")
	if err != nil {
		t.FailNow()
	}

	res := oanda.GetCandlesRes{}

	err = json.Unmarshal(contents, &res)
	if err != nil {
		t.FailNow()
	}

	events := oanda.AnalyzeCandles(res.Candles)
	_ = events

	chart := charts.NewKLine()
	chart.SetGlobalOptions(charts.WithTitleOpts(opts.Title{
		Title: "USD JPY S30",
	}))

	xs := []time.Time{}
	cs := []opts.KlineData{}

	loc, err := time.LoadLocation("Australia/Brisbane")
	if err != nil {
		t.FailNow()
	}

	for _, c := range res.Candles {
		xs = append(xs, c.Time.In(loc))
		cs = append(cs, opts.KlineData{
			Value: []float64{c.Mid.O, c.Mid.C, c.Mid.H, c.Mid.L},
		})
	}

	// for _, e := range events {
	// 	label := ""
	// 	for _, i := range e.Events {
	// 		label += i.Name + " (" + i.Value + ") " + ", "
	// 	}
	// 	cs[int(e.Index)].Name = label
	// }

	chart.SetXAxis(xs).AddSeries("prices", cs).
		SetSeriesOptions(
			charts.WithMarkPointNameTypeItemOpts(opts.MarkPointNameTypeItem{
				Name:     "highest value",
				Type:     "max",
				ValueDim: "highest",
			}),
			charts.WithMarkPointNameTypeItemOpts(opts.MarkPointNameTypeItem{
				Name:     "lowest value",
				Type:     "min",
				ValueDim: "lowest",
			}),
			charts.WithMarkPointStyleOpts(opts.MarkPointStyle{
				Label: &opts.Label{
					Show: opts.Bool(true),
				},
			}),
			charts.WithItemStyleOpts(opts.ItemStyle{
				Color:        "#ec0000",
				Color0:       "#00da3c",
				BorderColor:  "#8A0000",
				BorderColor0: "#008F28",
			}),
		)

	chart.SetGlobalOptions(
		charts.WithInitializationOpts(opts.Initialization{
			Width:  "1800px",
			Height: "1000px",
		}),
	)

	chart.SetGlobalOptions(
		charts.WithYAxisOpts(opts.YAxis{
			Scale: opts.Bool(true),
		}),
		charts.WithDataZoomOpts(opts.DataZoom{
			Start:      50,
			End:        100,
			XAxisIndex: []int{0},
		}),
	)

	f, _ := os.Create("kline.html")
	chart.Render(f)
}
