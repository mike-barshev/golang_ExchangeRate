package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
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

	// Postgres

	connStr := "host=localhost port=5432 user=postgres password=changeme dbname=postgres sslmode=disable"
	db, err := sqlx.Connect("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	// GET json
	var dailyGET func()
	dailyGET = func() {
		for {
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

				// Write in to DB

				sqlStatement := `
				INSERT INTO data ("Timestamp", "EUR", "USD")
				VALUES ($1, $2, $3)`
				_, err = db.Exec(sqlStatement, dataExRates.Timestamp, dataExRates.Valutes["EUR"].Value, dataExRates.Valutes["USD"].Value)
				if err != nil {
					log.Fatal(err)
				}
			}
			fmt.Println(time.Now())
			time.Sleep(24 * time.Hour)
			dailyGET()
		}
	}
	go dailyGET()

	// Server start

	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	e.Logger.Fatal(e.Start(":1323"))
}
