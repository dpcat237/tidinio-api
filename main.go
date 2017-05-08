package main

import (
	"github.com/tidinio/app"
	"log"
	"net/http"
)

func main() {
	app.InitializeRequiredData()

	router := app.NewRouter()
	log.Fatal(http.ListenAndServe(":80", router))
}
