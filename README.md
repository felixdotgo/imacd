# imacd
A Golang implementation of [LazyBear's IMACD](https://www.tradingview.com/v/qt6xLfLi/) indicator on TradingView

## Disclaimer
This is just a port of the indicator from TradingView to Golang. I made it as a learning exercise and not as a production ready indicator. Use at your own risk.

## Example
Query price candles from Binance API for BTCUSDT with 1h interval for the last 31 days.
```go
package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/felixdotgo/imacd"

	"github.com/fatih/color"
)

type Candle struct {
	OpenTime      float64 `json:"open_time"`
	Open          string  `json:"open"`
	High          string  `json:"high"`
	Low           string  `json:"low"`
	Close         string  `json:"close"`
	Volume        string  `json:"volume"`
	CloseTime     float64 `json:"close_time"`
	Qav           string  `json:"qav"`
	NumTrades     float64 `json:"num_trades"`
	TakerBaseVol  string  `json:"taker_base_vol"`
	TakerQuoteVol string  `json:"taker_quote_vol"`
	Ignore        string  `json:"ignore"`
}

func GetPriceCandles(symbol string, interval string, startTime int64, endTime int64) ([]Candle, error) {
	url := "https://api.binance.com/api/v3/klines"
	params := fmt.Sprintf("symbol=%s&interval=%s&startTime=%d&endTime=%d&limit=1000",
		symbol, interval, startTime, endTime)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.URL.RawQuery = params

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Unmarshal the response
	var candles [][]interface{}
	err = json.Unmarshal(body, &candles)
	if err != nil {
		return nil, err
	}

	// Convert the candles to the desired format
	var result []Candle
	for _, candle := range candles {
		result = append(result, Candle{
			OpenTime:      candle[0].(float64),
			Open:          candle[1].(string),
			High:          candle[2].(string),
			Low:           candle[3].(string),
			Close:         candle[4].(string),
			Volume:        candle[5].(string),
			CloseTime:     candle[6].(float64),
			Qav:           candle[7].(string),
			NumTrades:     candle[8].(float64),
			TakerBaseVol:  candle[9].(string),
			TakerQuoteVol: candle[10].(string),
			Ignore:        candle[11].(string),
		})
	}

	return result, nil
}

func main() {
	symbol := "BTCUSDT"
	interval := "1h"
	startTime := time.Now().Add(-24*31*time.Hour).UTC().UnixNano() / int64(time.Millisecond)
	endTime := time.Now().UTC().UnixNano() / int64(time.Millisecond)
	indicator := imacd.NewDefaultImpulseMACD()

	candles, err := GetPriceCandles(symbol, interval, startTime, endTime)
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, candle := range candles {
		openTime, _ := strconv.Atoi(fmt.Sprintf("%0.f", candle.OpenTime))
		closeTime, _ := strconv.Atoi(fmt.Sprintf("%0.f", candle.CloseTime))

		high, _ := strconv.ParseFloat(candle.High, 64)
		low, _ := strconv.ParseFloat(candle.Low, 64)
		close, _ := strconv.ParseFloat(candle.Close, 64)

		result := indicator.Update(high, low, close)

		color.New(color.Bold).Printf("OpenTime: ")
		fmt.Printf("%s, ", time.Unix(int64(openTime)/1000, 0).UTC().Format(time.RFC822Z))

		color.New(color.Bold).Printf("CloseTime: ")
		fmt.Printf("%s, ", time.Unix(int64(closeTime)/1000, 0).UTC().Format(time.RFC822Z))


		color.New(color.Bold, color.FgBlue).Printf("MD: ")
		fmt.Printf("%.f, ", result.MD)

		color.New(color.Bold, color.FgHiYellow).Printf("Signal: ")
		fmt.Printf("%.f, ", result.SB)

		color.New(color.Bold).Printf("Histogram: ")
		fmt.Printf("%.f, ", result.SH)

		fmt.Printf("%s\n", result.Color)
	}
}
```

## License
MIT License

Copyright (c) 2025 xileF

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
