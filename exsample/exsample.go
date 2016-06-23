package main

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/ieee0824/louis"
)

func main() {
	louis.AddClient(&http.Client{})
	c := louis.NewClient()
	req, _ := http.NewRequest("GET", "https://www.cman.jp/network/support/go_access.cgi", nil)
	r, _ := c.Do(req)
	bin, _ := ioutil.ReadAll(r.Body)
	fmt.Println(string(bin))
}
