package oneid

import (
    broker "github.com/longguikeji/arkid-broker/base_broker"
    handler "github.com/longguikeji/arkid-broker/handler"

    "github.com/gorilla/mux"
)

// GetRouter return the router
func GetRouter() *mux.Router {
    var r = mux.NewRouter()

    r.HandleFunc("/siteapi/v1/user/", broker.ParseBroker(handler.UserListBroker{}))
    r.HandleFunc("/siteapi/v1/user/{username:[0-9a-zA-Z]+}/", broker.ParseBroker(handler.UserDetailBroker{}))
    r.HandleFunc("/siteapi/v1/ucenter/login/", broker.ParseBroker(handler.UserLoginBroker{}))
    r.HandleFunc("/siteapi/v1/auth/token/", broker.ParseBroker(handler.AuthTokenBroker{}))

    r.PathPrefix("/siteapi/v1").HandlerFunc(broker.ParseBroker(broker.DefaultBroker{}))

    return r
}
