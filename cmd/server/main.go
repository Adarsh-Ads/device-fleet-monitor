package main

import (
	"fmt"
	"net/http"
)

func home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Device Fleet Monitor")
}

func main() {
	http.HandleFunc("/", home)

	fmt.Println("Starting Device Fleet Monitor on :8080")

	err := http.ListenAndServe(":8080", nil)

	if err != nil {
		fmt.Println("Server failed:", err)
	}
}
