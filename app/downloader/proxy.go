package downloader

import (
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"strings"

	"golang.org/x/net/proxy"
)

func LoadProxy(file string) ([]*url.URL, error) {
	b, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	var proxys []*url.URL

	strs := strings.Split(string(b), "\n")
	for _, str := range strs {
		if proxy, err := url.Parse(str); err == nil {
			proxys = append(proxys, proxy)
		} else {
			log.Println("invalid proxy: ", err)
		}
	}
	return proxys, nil
}

func NewClientWithProxy(uri *url.URL) *http.Client {
	tr := &http.Transport{}
	if dialer, err := proxy.FromURL(uri, &net.Dialer{}); err == nil {
		tr.Dial = dialer.Dial
	} else {
		tr.Proxy = http.ProxyURL(uri)
	}
	return &http.Client{
		Transport: tr,
	}
}
