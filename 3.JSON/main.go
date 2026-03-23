package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
)

type Payment struct {
	Description string `json:"description"`
	USD         int    `json:"usd"`
	FullName    string `json:"fullname"`
	Address     string `json:"address"`
}
type HTTPResponse struct {
	Money          int
	PaymentHistory []Payment
}

var balance = 1000
var paymentHistory = make([]Payment, 0)
var mtx = sync.Mutex{}

func (p Payment) Println() {
	fmt.Println("Description -", p.Description)
	fmt.Println("USD -", p.USD)
	fmt.Println("Full Name -", p.FullName)
	fmt.Println("Address -", p.Address)
}

func payHandler(w http.ResponseWriter, r *http.Request) {
	var payment Payment

	if err := json.NewDecoder(r.Body).Decode(&payment); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// httpBody, err := io.ReadAll(r.Body)
	// if err != nil {
	// 	w.WriteHeader(http.StatusBadRequest)
	// 	return
	// }

	// if err := json.Unmarshal(httpBody, &payment); err != nil {
	// 	fmt.Println("Произошла ошибка:", err)
	// 	w.WriteHeader(http.StatusInternalServerError)
	// 	return
	// }

	mtx.Lock()
	if balance-payment.USD > 0 {
		balance -= payment.USD
	}

	paymentHistory = append(paymentHistory, payment)

	httpResponse := HTTPResponse{
		Money:          balance,
		PaymentHistory: paymentHistory,
	}

	b, err := json.MarshalIndent(httpResponse, "", "	")
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if _, err := w.Write(b); err != nil {
		fmt.Println(err)
		return
	}

	payment.Println()

	mtx.Unlock()

}

func main() {
	http.HandleFunc("/pay", payHandler)

	if err := http.ListenAndServe(":9091", nil); err != nil {
		fmt.Println("Произошла ошибка сервера:", err)
	}

}
