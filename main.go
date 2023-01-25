package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {
	getReq("https://google.com")
}

func getReq(url string) {
	response, err := http.Get(url)
	if err != nil {
		log.Fatalln(err)
	}
	defer response.Body.Close()

	rawBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(string(rawBody))
}
