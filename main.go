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
	MaxDur := int(dur.Milliseconds())
	TX, err := ConvertIn(tx)
	matr := make([][]float64, len(tx)+1)
	for i := 0; i < len(tx)+1; i++ {
		matr[i] = make([]float64, MaxDur+1)
	}

	for i := 0; i <= len(TX); i++ {
		for j := 0; j <= MaxDur; j++ {
			//fmt.Println(i, j)
			if i == 0 || j == 0 {
				matr[i][j] = 0
			} else {
				if TX[i-1].API > j {
					matr[i][j] = matr[i-1][j]
				} else {
					prev := matr[i-1][j]
					byFormula := TX[i-1].Amount + matr[i-1][j-TX[i-1].API]
					if prev > byFormula {
						if i == len(TX) {
							ans = append(ans, tx[i-2])
						}
						matr[i][j] = prev
					} else {
						if i == len(TX) {
							ans = append(ans, tx[i-1])
						}
						matr[i][j] = byFormula
					}
				}
			}
		}
	}
	fmt.Println(matr[len(TX)][MaxDur])

	return
}

func main() {
	fmt.Println(time.Now())

	transactions := CsvParse("transactions.csv")

	trans, _ := prioritize(transactions, 1*time.Second)
	sum := .0
	for i := 0; i < len(trans); i++ {
		tmp, _ := strconv.ParseFloat(trans[i].Amount, 64)
		sum += tmp
	}
	fmt.Println(sum)
	fmt.Println(time.Now())

}
