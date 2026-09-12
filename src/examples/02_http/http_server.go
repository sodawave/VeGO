package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "hello vego")
	})
	fmt.Println("listening on :8088")
	_ = http.ListenAndServe(":8088", nil)
}
