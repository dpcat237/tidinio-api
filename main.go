package main

/**
More info:
https://github.com/gorilla/mux
https://github.com/avelino/awesome-go
http://jinzhu.me/gorm/crud.html#query
*/

/**
TODO:
- API: user register, login
- API: get user new articles
*/

import (
	"github.com/tidinio/app"
	"log"
	"net/http"
)

func main() {
	router := app.NewRouter()

	log.Fatal(http.ListenAndServe(":80", router))
}
