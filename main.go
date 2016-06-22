package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"

	"golang.org/x/net/proxy"

	"github.com/PuerkitoBio/goquery"
)

func init() {
	//addClient()
	addClient()
	clients = append(clients, newTorClient("localhost", "9050"))
	clients = append(clients, newTorClient("localhost", "9052"))
}

var clients = []*http.Client{}
var clientNum = 0

func newTorClient(host string, port string) *http.Client {
	tbProxyURL, err := url.Parse("socks5://" + host + ":" + port)
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

func addClient() {
	clients = append(clients, &http.Client{})
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

func client(req <-chan *http.Request, resp chan<- *http.Response) {
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

func main() {
	req := make(chan *http.Request, 256)
	resp := make(chan *http.Response, 256)
	client(req, resp)

	r0, err := http.NewRequest("GET", "https://www.cman.jp/network/support/go_access.cgi", nil)
	if err != nil {
		log.Fatalln(err)
	}
	req <- r0

	r1, err := http.NewRequest("GET", "https://www.cman.jp/network/support/go_access.cgi", nil)
	if err != nil {
		log.Fatalln(err)
	}
	req <- r1

	r2, err := http.NewRequest("GET", "https://www.cman.jp/network/support/go_access.cgi", nil)
	if err != nil {
		log.Fatalln(err)
	}
	req <- r2

	for {
		select {
		case res := <-resp:
			doc, _ := goquery.NewDocumentFromReader(res.Body)
			fmt.Println(doc.Find(".outIp").Text())
		default:

		}
	}

}
