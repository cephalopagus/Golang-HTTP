package main

import (
	"fmt"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	str := "hello world"
	arr := []byte(str)

	_, err := w.Write(arr)

	if err != nil {
		fmt.Println("Произошла ошибка:", err.Error())
	} else {
		fmt.Println("Я корректно обработал HTTP запрос!")
	}
}
func cancelHandler(w http.ResponseWriter, r *http.Request) {
	str := "Оплата отменета!"
	arr := []byte(str)

	_, err := w.Write(arr)

	if err != nil {
		fmt.Println("Произошла ошибка:", err.Error())
	} else {
		fmt.Println("Я корректно отменил оплату!")
	}
}
func payHandler(w http.ResponseWriter, r *http.Request) {
	str := "Оплата произведена!"
	arr := []byte(str)

	_, err := w.Write(arr)

	if err != nil {
		fmt.Println("Произошла ошибка:", err.Error())
	} else {
		fmt.Println("Я корректно произвел оплату!")
	}
}

func main() {
	http.HandleFunc("/default", handler)
	http.HandleFunc("/pay", payHandler)
	http.HandleFunc("/cancel", cancelHandler)

	fmt.Println("Запускаю сервер")
	err := http.ListenAndServe(":9091", nil)
	if err != nil {
		fmt.Println("Произошла ошибка сервера:", err.Error())
	}
}
