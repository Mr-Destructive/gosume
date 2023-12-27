package main

import (
	"fmt"
	"log"
	"net/http"
    "github.com/mr-destructive/gosume/api"
)

func handleRequests(port int) {
	http.HandleFunc("/", api.Handler)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}

func main() {
	port := 8080
	log.Println(port)
	handleRequests(port)
}
