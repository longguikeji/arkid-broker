package oneid

import (
    "net/http"

    broker "github.com/longguikeji/arkid-broker/base_broker"
)

// UserLoginBroker ...
type UserLoginBroker struct {
    broker.Broker
}

// ProcessRequest ...
func (b UserLoginBroker) ProcessRequest(r *http.Request) {
}

// ProcessResponse ...
func (b UserLoginBroker) ProcessResponse(r *http.Response) error {
    res, err := b.InitResponse(r)
    if err != nil {
        return err
    }

    switch r.StatusCode {
    case 200:
        res.Status.Message = "登录成功"
    case 400:
        res.Status.Message = "用户名密码错误"
    }
    return b.WrapResponse(res, r)
}

// AuthTokenBroker ...
type AuthTokenBroker struct {
    broker.Broker
}

// ProcessRequest ...
func (b AuthTokenBroker) ProcessRequest(r *http.Request) {
}

// ProcessResponse ...
func (b AuthTokenBroker) ProcessResponse(r *http.Response) error {
    res, err := b.InitResponse(r)
    if err != nil {
        return err
    }

    switch r.StatusCode {
    case 200:
        res.Status.Message = "认证成功"
    case 401:
        res.Status.Message = "Token无效或未提供"
    case 403:
        res.Status.Message = "该Token权限不足"
    }
    return b.WrapResponse(res, r)
}
