package main

import (
	"fmt"
	"net/http"

	"github.com/ieee0824/louis"
)

func main() {
	louis.AddClient(&http.Client{})

	fmt.Println(louis.NumClient())
	fmt.Println(louis.ClientIPList())
	fmt.Println(louis.ClientIPList())
}
