package main

import (
	"encoding/json"
	"fmt"
	"github.com/gocarina/gocsv"
	"io/ioutil"
	"math"
	"os"
	"strconv"
	"time"
)

type (
	//Trans структура для обработанных транзакций
	Trans struct {
		ID              string
		Amount          float64
		BankName        string
		BankCountryCode string
		Time            int //Время прихода ответа от банка
	}

	Transaction struct {
		ID              string `csv:"id"`
		Amount          string `csv:"amount"`
		BankName        string
		BankCountryCode string `csv:"bank_country_code"`
	}
)

const defaultMaxTxDuration = time.Second

//Convert Обрабатываем slice []Transaction
func Convert(tx []Transaction) (TX []Trans, err error) {
	Json, err := os.Open("api_latencies.json")
	if err != nil {
		return
	}
	defer Json.Close()

	byteValue, _ := ioutil.ReadAll(Json)

	var api map[string]int
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
		tmp.Time = api[i.BankCountryCode]
		TX = append(TX, tmp)
	}
	return
}

//CsvParse Из Csv файла с транзакциями делаем slice []Transaction
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

//prioritize Ищем подходящие транзакции алгоритмом динамического программирования
func prioritize(tx []Transaction, dur time.Duration) (ans []Transaction, err error) {
	if dur == 0 {
		dur = defaultMaxTxDuration
	}
	MaxDur := int(dur.Milliseconds())
	TX, err := Convert(tx)
	matr := make([][]float64, len(tx)+1)
	for i := 0; i < len(tx)+1; i++ {
		matr[i] = make([]float64, MaxDur+1)
	}

	for i := 0; i <= len(TX); i++ { //Ищем максимальное количество Долларов
		for j := 0; j <= MaxDur; j++ {
			if i == 0 || j == 0 {
				matr[i][j] = 0
			} else {
				if TX[i-1].Time > j {
					matr[i][j] = matr[i-1][j]
				} else {
					prev := matr[i-1][j]
					Formula := TX[i-1].Amount + matr[i-1][j-TX[i-1].Time]
					matr[i][j] = math.Max(prev, Formula)
				}
			}
		}
	}

	for i := len(TX); i > 0; i-- { //Составляем slice нужных транзакций
		if matr[i][MaxDur] > matr[i-1][MaxDur] {
			ans = append(ans, tx[i-1])
			MaxDur -= TX[i-1].Time
		}
	}
	return
}

func GetSum(tx []Transaction) (sum float64) { //Получаем сумму транзакций из slice
	for i := 0; i < len(tx); i++ {
		tmp, _ := strconv.ParseFloat(tx[i].Amount, 64)
		sum += tmp
	}
	return
}

func main() {
	transactions := CsvParse("transactions.csv")

	TimeNow := time.Now()
	trans, _ := prioritize(transactions, 1000*time.Millisecond)
	fmt.Printf("Answer question 1 \n(1 second): %.2f(USD), tiwe work: %v(ms)\n\n", GetSum(trans), time.Now().Sub(TimeNow).Milliseconds())

	fmt.Println("Answer question 2 :")
	TimeNow = time.Now()
	trans, _ = prioritize(transactions, 50*time.Millisecond)
	fmt.Printf("(50ms): %.2f(USD), tiwe work: %v(ms)\n", GetSum(trans), time.Now().Sub(TimeNow).Milliseconds())
	TimeNow = time.Now()
	trans, _ = prioritize(transactions, 60*time.Millisecond)
	fmt.Printf("(60ms): %.2f(USD), tiwe work: %v(ms)\n", GetSum(trans), time.Now().Sub(TimeNow).Milliseconds())
	TimeNow = time.Now()
	trans, _ = prioritize(transactions, 90*time.Millisecond)
	fmt.Printf("(90ms): %.2f(USD), tiwe work: %v(ms)\n", GetSum(trans), time.Now().Sub(TimeNow).Milliseconds())

}
