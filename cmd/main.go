package main

import (
	"fmt"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello World"))
}

func main() {
	http.HandleFunc("/", handler)

	fmt.Println("Escutando na porta :8080...")
	err := http.ListenAndServe(":8080", nil)

	if err != nil {
		panic(err)
	}
}
