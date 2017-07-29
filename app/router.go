package app

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/tidinio/src/controller/api"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

var routes = Routes{
	/** Article **/
	Route{
		"Add shared article",
		"POST",
		"/article/add_shared",
		api_controller.AddSharedItem,
	},
	Route{
		"Sync items",
		"POST",
		"/article/sync",
		api_controller.SyncItems,
	},
	Route{
		"Get saved articles",
		"POST",
		"/saved_article/list",
		api_controller.ListTagItems,
	},
	Route{
		"Sync saved articles changes",
		"POST",
		"/saved_article/sync",
		api_controller.SyncTagItems,
	},
	/** Device **/
	Route{
		"Update push notification ID for device. Now Firebase",
		"PUT",
		"/device/notification",
		api_controller.UpdateNotificationId,
	},
	/** Feed **/
	Route{
		"Add feed, subscribe user to this feed and add last items for user",
		"POST",
		"/feed/add",
		api_controller.AddFeed,
	},
	Route{
		"Unsubscribe from feed",
		"DELETE",
		"/feed/{id}",
		api_controller.DeleteFeed,
	},
	Route{
		"Edit feed title",
		"PATCH",
		"/feed/{id}",
		api_controller.EditFeed,
	},
	Route{
		"Get feed sources",
		"GET",
		"/feed/sources",
		api_controller.GetSources,
	},
	Route{
		"Sync feeds",
		"POST",
		"/feed/sync",
		api_controller.SyncFeeds,
	},
	/** Filter **/
	Route{
		"Add filters",
		"POST",
		"/filter",
		api_controller.AddFilters,
	},
	Route{
		"Delete filters",
		"DELETE",
		"/filter",
		api_controller.DeleteFilters,
	},
	Route{
		"Get filters",
		"GET",
		"/filter",
		api_controller.GetFilters,
	},
	Route{
		"Update filters",
		"PUT",
		"/filter",
		api_controller.UpdateFilters,
	},
	/** Tag **/
	Route{
		"Add tags",
		"POST",
		"/tag",
		api_controller.AddTags,
	},
	Route{
		"Delete tags",
		"DELETE",
		"/tag",
		api_controller.DeleteTags,
	},
	Route{
		"Get tags",
		"GET",
		"/tag",
		api_controller.GetTags,
	},
	Route{
		"Update tags",
		"PUT",
		"/tag",
		api_controller.UpdateTags,
	},
	/** User **/
	Route{
		"User feedback",
		"POST",
		"/user/feedback",
		api_controller.UserFeedback,
	},
	Route{
		"User login",
		"POST",
		"/user/login",
		api_controller.UserLogin,
	},
	Route{
		"User register",
		"POST",
		"/user/register",
		api_controller.UserRegister,
	},
	Route{
		"User password recovery",
		"PUT",
		"/user/password_request/{hash}",
		api_controller.UserPasswordRecovery,
	},
	Route{
		"User password recovery request",
		"POST",
		"/user/password_request",
		api_controller.UserPasswordRecoveryRequest,
	},
	Route{
		"User preview login",
		"POST",
		"/user/preview/login",
		api_controller.UserPreviewLogin,
	},
	Route{
		"User preview register",
		"POST",
		"/user/preview/register",
		api_controller.UserPreviewRegister,
	},
	Route{
		"User social login",
		"POST",
		"/user/social/login",
		api_controller.UserSocialLogin,
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
