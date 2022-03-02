package main

import (
	"encoding/json"
	"fmt"
	"github.com/gocarina/gocsv"
	"io/ioutil"
	"os"
	"strconv"
	"time"
)

type (
	Trans struct {
		ID              string
		Amount          float64
		BankName        string
		BankCountryCode string
		API             int
	}

	Transaction struct {
		ID              string `csv:"id"`
		Amount          string `csv:"amount"`
		BankName        string
		BankCountryCode string `csv:"bank_country_code"`
	}
)

const defaultMaxTxDuration = time.Second

func ConvertIn(tx []Transaction) (TX []Trans, err error) {
	Json, err := os.Open("api_latencies.json")
	if err != nil {
		return
	}
	defer Json.Close()

	byteValue, _ := ioutil.ReadAll(Json)

	var api map[string]string
	json.Unmarshal([]byte(byteValue), &api)

	TX = make(
		[]Trans,
		0,
		len(tx),
	)
	for _, i := range tx {
		var tmp Trans
		tmp.ID = i.ID
		tmp.Amount, _ = strconv.ParseFloat(i.Amount, 64)
		tmp.BankName = i.BankName
		tmp.BankCountryCode = i.BankCountryCode
		tmp.API, _ = strconv.Atoi(api[i.BankCountryCode])
		TX = append(TX, tmp)
	}
	return
}

func ConvertOut(tx []Trans) (TX []Transaction) {
	TX = make(
		[]Transaction,
		0,
		len(tx),
	)
	for _, i := range tx {
		var tmp Transaction
		tmp.ID = i.ID
		tmp.Amount = strconv.FormatFloat(i.Amount, 'E', -1, 64)
		tmp.BankName = i.BankName
		tmp.BankCountryCode = i.BankCountryCode
		TX = append(TX, tmp)
	}
	return
}
func CsvParse(file string) (transactions []Transaction) {
	fmt.Println(time.Now())
	db, err := os.Open(file)
	defer db.Close()

	if err != nil {
		panic(err)
	}
	if err = gocsv.UnmarshalFile(db, &transactions); err != nil {
		panic(err)
	}
	return
}
func prioritize(tx []Transaction, dur time.Duration) (ans []Transaction, err error) {
	if dur == 0 {
		dur = defaultMaxTxDuration
	}
	TX, err := ConvertIn(tx)
	var matr [len(TX)][]float64
	return
}

func main() {
	fmt.Println(time.Now())

	transactions := CsvParse("transactions.csv")

	prioritize(transactions, 1*time.Second)

	fmt.Println(time.Now())

}
