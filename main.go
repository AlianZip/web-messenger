package main

import (
	"fmt"
	"log"
	"messenger/routes"
	"net/http"
)

func main() {
	r := routes.NewRouter()

	fmt.Println("http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
