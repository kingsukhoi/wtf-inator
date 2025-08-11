package proxy

import (
	"net/http/httputil"
	"net/url"
)

type WtfProxy struct {
	Proxy *httputil.ReverseProxy
}

func NewWtfProxy(host string) (*WtfProxy, error) {
	currUrl, err := url.Parse(host)
	if err != nil {
		return nil, err
	}

	rtnMe := &WtfProxy{
		Proxy: httputil.NewSingleHostReverseProxy(currUrl),
	}

	rtnMe.Proxy.ModifyResponse = responseHandler

	return rtnMe, nil
}
