package broker

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
)

// IBroker ...
type IBroker interface {
	ProcessRequest(*http.Request)
	ProcessResponse(*http.Response) error
	LogResponse(*http.Response)
	LogRequest(*http.Request)
}

// Broker ...
type Broker struct {
}

// DefaultBroker ...
type DefaultBroker struct {
	Broker
}

// ProcessRequest of DefaultBroker
func (b DefaultBroker) ProcessRequest(r *http.Request) {
}

// ProcessResponse of DefaultBroker
func (b DefaultBroker) ProcessResponse(r *http.Response) error {
	res, err := b.InitResponse(r)
	if err != nil {
		return err
	}
	return b.WrapResponse(res, r)
}

// ResStatus ...
type ResStatus struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// ResContent ...
type ResContent struct {
	Status ResStatus   `json:"status"`
	Result interface{} `json:"result"`
}

// AcceptJSON 设置请求头 Accept: application/json
func (b *Broker) AcceptJSON(r *http.Request) {
	r.Header.Set("Accept", "application/json")
}

// ReadResponse 读取response 内容
func (b *Broker) ReadResponse(c *ResContent, r *http.Response) (err error) {
	var reader io.ReadCloser
	switch r.Header.Get("Content-Encoding") {
	case "gzip":
		reader, _ = gzip.NewReader(r.Body)
		defer reader.Close()
		r.Header.Del("Content-Encoding")
	default:
		reader = r.Body
	}
	body, err := ioutil.ReadAll(reader)
	if err != nil {
		return err
	}
	err = json.Unmarshal(body, &(c.Result))
	if err != nil {
		c.Result = string(body)
	}
	return nil
}

// InitResponse 处理常见标准错误并停供默认提示
func (b *Broker) InitResponse(r *http.Response) (*ResContent, error) {
	res := &ResContent{
		Status: ResStatus{
			Code:    999,
			Message: "请求失败(未知错误)",
		},
	}
	err := b.ReadResponse(res, r)
	if err != nil {
		return res, err
	}
	switch r.StatusCode {
	case 200:
		res.Status.Code = 0
		res.Status.Message = "请求成功"
	case 201:
		res.Status.Code = 0
		res.Status.Message = "创建成功"
	case 204:
		res.Status.Code = 0
		res.Status.Message = "删除成功"
	case 400:
		res.Status.Code = 400
		res.Status.Message = "请求参数错误"
	case 401:
		res.Status.Code = 401
		res.Status.Message = "无效的身份凭证"
	case 403:
		res.Status.Code = 403
		res.Status.Message = "权限不足"
	case 404:
		res.Status.Code = 404
		res.Status.Message = "该记录不存在"
	case 405:
		res.Status.Code = 405
		res.Status.Message = "接口使用方式错误"
	}

	if r.StatusCode >= 500 {
		res.Status.Code = 500
		res.Status.Message = "接口内部异常"
	}
	return res, nil
}

// WrapResponse 将ResContent以json返回
func (b *Broker) WrapResponse(c *ResContent, r *http.Response) error {
	status := map[string]interface{}{
		"code":    c.Status.Code,
		"message": c.Status.Message,
	}
	res := map[string]interface{}{
		"status": status,
		"result": c.Result,
	}
	content, err := json.Marshal(&res)
	if err != nil {
		return err
	}

	r.Body = ioutil.NopCloser(bytes.NewReader(content))
	r.Header.Set("Content-Type", "application/json")
	r.Header.Set("Content-Length", strconv.Itoa(len(content)))
	return nil
}

// LogResponse ...
func (b Broker) LogResponse(r *http.Response) {
	// TODO
	fmt.Println(r.Body)
}

// LogRequest ...
func (b Broker) LogRequest(r *http.Request) {
	// TODO
	fmt.Println(r.Body)
}
