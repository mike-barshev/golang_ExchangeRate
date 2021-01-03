package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

// ExchangeRates - http get json struct
type ExchangeRates struct {
	Date         time.Time         `json:"Date"`
	PreviousDate time.Time         `json:"PreviousDate"`
	PreviousURL  string            `json:"PreviousURL"`
	Timestamp    time.Time         `json:"Timestamp"`
	Valutes      map[string]Valute `json:"Valute"`
}

// Valute - valute struct
type Valute struct {
	ID       string  `json:"ID"`
	NumCode  string  `json:"NumCode"`
	CharCode string  `json:"CharCode"`
	Nominal  int     `json:"Nominal"`
	Name     string  `json:"Name"`
	Value    float64 `json:"Value"`
	Previous float64 `json:"Previous"`
}

func main() {

	response, err := http.Get("https://www.cbr-xml-daily.ru/daily_json.js")
	if err != nil {
		fmt.Printf("HTTP request failed, error: %s\n", err)
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		dataExRates := ExchangeRates{}
		jsonErr := json.Unmarshal(data, &dataExRates)
		if jsonErr != nil {
			fmt.Printf("Json unmarshal error: %s\n", jsonErr)
		}
		fmt.Println(dataExRates)
		// fmt.Println(dataExRates.Valutes["USD"])
	}
}
