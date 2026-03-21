package main

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"sync"
	"sync/atomic"
)

var money = atomic.Int64{}
var bank = atomic.Int64{}
var mtx = sync.Mutex{}

func payHandler(w http.ResponseWriter, r *http.Request) {
	httpRequestBody, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println("failed reading from http body:", err)
		return
	}
	httpRequest := string(httpRequestBody)
	paymentAmount, err := strconv.Atoi(httpRequest)
	if err != nil {
		fmt.Println("failed convert str into int:", err)
		return
	}
	mtx.Lock()
	if money.Load()-int64(paymentAmount) >= 0 {
		money.Add(-int64(paymentAmount))
		fmt.Println("Оплата прошла успешно:", money.Load())
	} else {
		fmt.Println("Не хвататет денег для проведения оплаты!")
	}
	mtx.Unlock()

}

func saveHandler(w http.ResponseWriter, r *http.Request) {
	httpRequestBody, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println("failed reading from http body:", err)
		return
	}

	httpRequest := string(httpRequestBody)

	saveAmount, err := strconv.Atoi(httpRequest)
	if err != nil {
		fmt.Println("failed convert str into int:", err)
		return
	}

	mtx.Lock()
	if money.Load() >= int64(saveAmount) {

		money.Add(-int64(saveAmount))
		bank.Add(int64(saveAmount))

		fmt.Println("Новое значение переменной money:", money.Load())

		fmt.Println("Новое значение переменной bank:", bank.Load())
	} else {
		fmt.Println("Не хвататет денег чтобы положить в копилку!")
	}
	mtx.Unlock()
}

func main() {
	money.Add(1000)
	http.HandleFunc("/pay", payHandler)
	http.HandleFunc("/save", saveHandler)

	err := http.ListenAndServe(":9091", nil)
	if err != nil {
		fmt.Println("Произошла ошибка сервера:", err.Error())
	}

}
