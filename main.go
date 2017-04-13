package main

import (
	"log"
	"net/http"
	"github.com/tidinio/app"
)

func main() {
	app.InitializeRequiredData()

	router := app.NewRouter()
	log.Fatal(http.ListenAndServe(":80", router))
}
