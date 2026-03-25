package main

import (
	"fmt"
	"study/http"
	"study/todo"
)

func main() {
	todoList := todo.NewList()
	httpHandlers := http.NewHTTPHandler(todoList)
	httpServer := http.NewHTTPServer(httpHandlers)
	if err := httpServer.StartServer(); err != nil {
		fmt.Println("err serv", err)
	}
}
