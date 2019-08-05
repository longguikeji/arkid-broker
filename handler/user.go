package oneid

import (
    "net/http"

    broker "github.com/rockl2e/ark-apisvr/broker"
)

// UserListBroker ...
type UserListBroker struct {
    broker.Broker
}

// ProcessRequest ...
func (b UserListBroker) ProcessRequest(r *http.Request) {
    r.Header.Set("ARKER", "oneid_broker")
}

// ProcessResponse ...
func (b UserListBroker) ProcessResponse(r *http.Response) error {
    res, err := b.InitResponse(r)
    if err != nil {
        return err
    }

    return b.WrapResponse(res, r)
}

// UserDetailBroker ...
type UserDetailBroker struct {
  broker.Broker
}

// ProcessRequest ...
func (b UserDetailBroker) ProcessRequest(r *http.Request) {
    r.Header.Set("ARKER", "oneid_broker")
}


// ProcessResponse ...
func (b UserDetailBroker) ProcessResponse(r *http.Response) error {
    res, err := b.InitResponse(r)
    if err != nil {
        return err
    }

    return b.WrapResponse(res, r)
}
