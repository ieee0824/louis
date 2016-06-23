package louis

import (
	"net/http"
	"net/url"

	"golang.org/x/net/proxy"

	"github.com/PuerkitoBio/goquery"
)

var clients = []*http.Client{}
var clientNum = 0
var clientCacheFlag = false
var ips []string

type Client struct {
}

func NewClient() *Client {
	return &Client{}
}

func (c *Client) Do(req *http.Request) (*http.Response, error) {
	return getClient().Do(req)
}

func NewTorClient(host string, port string) *http.Client {
	url := "socks5://" + host + ":" + port
	return NewProxyClient(url)
}

func NewProxyClient(u string) *http.Client {
	clientCacheFlag = false
	tbProxyURL, err := url.Parse(u)
	if err != nil {
		return nil
	}
	tbDialer, err := proxy.FromURL(tbProxyURL, proxy.Direct)
	if err != nil {
		return nil
	}
	tbTransport := &http.Transport{Dial: tbDialer.Dial}
	return &http.Client{Transport: tbTransport}
}

func AddClient(c *http.Client) {
	clients = append(clients, c)
}

func getClient() *http.Client {
	maxNum := len(clients)
	if maxNum == clientNum {
		clientNum = 1
		return clients[0]
	}
	clientNum++
	return clients[clientNum-1]
}

func ClientPipe(req <-chan *http.Request, resp chan<- *http.Response) {
	go func() {
		for {
			select {
			case r := <-req:
				c := getClient()
				res, _ := c.Do(r)
				resp <- res
			default:
			}
		}
	}()
}

func NumClient() int {
	return len(clients)
}

func ClientIPList() []string {
	if clientCacheFlag {
		return ips
	}
	ret := []string{}

	for _, c := range clients {
		req, err := http.NewRequest("GET", "https://www.cman.jp/network/support/go_access.cgi", nil)
		if err != nil {
			continue
		}
		resp, err := c.Do(req)
		if err != nil {
			continue
		}
		defer resp.Body.Close()
		doc, err := goquery.NewDocumentFromReader(resp.Body)
		if err != nil {
			continue
		}
		result := doc.Find(".outIp").Text()
		ret = append(ret, result)
	}

	clientCacheFlag = true
	ips = ret
	return ret
}
