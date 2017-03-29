package app

import (
	"net/http"
	"github.com/gorilla/mux"
	"github.com/tidinio/src/api/controller"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

var routes = Routes{
	/** Feed **/
	Route{
		"Index",
		"GET",
		"/",
		controller.Index,
	},
	/** Item **/
	Route{
		"Add shared article",
		"POST",
		"/article/add_shared",
		controller.AddSharedItem,
	},
	Route{
		"Sync items",
		"POST",
		"/article/sync",
		controller.SyncItems,
	},
	Route{
		"Get saved articles",
		"POST",
		"/saved_article/list",
		controller.ListTagItems,
	},
	Route{
		"Sync saved articles changes",
		"POST",
		"/saved_article/sync",
		controller.SyncTagItems,
	},
	/** User **/
	Route{
		"User register",
		"POST",
		"/user",
		controller.UserRegister,
	},
}

func NewRouter() *mux.Router {

	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		var handler http.Handler

		handler = route.HandlerFunc
		handler = Logger(handler, route.Name)

		router.
		Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}

	return router
}
