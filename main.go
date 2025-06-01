package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/AlianZip/web-messenger/routes"
)

func main() {
	r := routes.NewRouter()
	fs := http.FileServer(http.Dir("static"))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))

	fmt.Println("http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
