package main

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"sync"
)

var money = 1000
var bank = 0
var mtx = sync.Mutex{}

func payHandler(w http.ResponseWriter, r *http.Request) {
	httpRequestBody, err := io.ReadAll(r.Body)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		msg := "failed reading from http body:" + err.Error()
		fmt.Println(msg)
		_, err := w.Write([]byte(msg))
		if err != nil {
			fmt.Println("failed send request:", err)
		}
		return
	}

	paymentAmount, err := strconv.Atoi(string(httpRequestBody))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		msg := "failed convert str into int:" + err.Error()
		fmt.Println(msg)
		_, err := w.Write([]byte(msg))
		if err != nil {
			fmt.Println("failed send request:", err)
		}
		return
	}

	mtx.Lock()
	defer mtx.Unlock()

	if money-paymentAmount >= 0 {
		money -= paymentAmount
		msg := "Оплата прошла успешно:" + strconv.Itoa(money)
		fmt.Println(msg)
		_, err := w.Write([]byte(msg))
		if err != nil {
			fmt.Println("failed send request:", err)
		}
	} else {
		_, err := w.Write([]byte("Не хватает денег для проведения оплаты!"))
		if err != nil {

			fmt.Println("Не хватает денег для проведения оплаты!")
		}
	}
}

func saveHandler(w http.ResponseWriter, r *http.Request) {
	httpRequestBody, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println("failed reading from http body:", err)
		return
	}

	saveAmount, err := strconv.Atoi(string(httpRequestBody))
	if err != nil {
		fmt.Println("failed convert str into int:", err)
		return
	}

	mtx.Lock()
	defer mtx.Unlock()

	if money >= saveAmount {
		money -= saveAmount
		bank += saveAmount

		fmt.Println("Новое значение переменной money:", money)
		fmt.Println("Новое значение переменной bank:", bank)
	} else {
		fmt.Println("Не хватает денег чтобы положить в копилку!")
	}
}

func main() {

	http.HandleFunc("/pay", payHandler)
	http.HandleFunc("/save", saveHandler)

	err := http.ListenAndServe(":9091", nil)
	if err != nil {
		fmt.Println("Произошла ошибка сервера:", err.Error())
	}
}
