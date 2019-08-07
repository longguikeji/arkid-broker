package broker

import (
	"crypto/tls"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/spf13/viper"
)

// ParseBroker convert broker to http handler
func ParseBroker(b IBroker) func(http.ResponseWriter, *http.Request) {
	wrapProcessRequest := func(req *http.Request) {
		upstreamURI := viper.GetString("target")
		upstream, _ := url.Parse(upstreamURI)

		req.URL.Host = upstream.Host
		req.URL.Scheme = upstream.Scheme
		req.Host = upstream.Host
		req.Header.Set("Accept", "application/json")
		b.ProcessRequest(req)
	}

	wrapProcessResponse := func(res *http.Response) error {
		b.LogResponse(res)
		err := b.ProcessResponse(res)
		b.LogResponse(res)
		return err
	}
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	proxy := httputil.ReverseProxy{Transport: tr}
	proxy.Director = wrapProcessRequest
	proxy.ModifyResponse = wrapProcessResponse
	return proxy.ServeHTTP
}
