package main

import (
	"log"
	"net/http"
	"os"
)

func main()  {
	SetupLogger()
	simpleHttpGet("www.google.com")
	simpleHttpGet("http://www.google.com")
}

func SetupLogger()  {
	file, _ := os.OpenFile("C:/Users/诚/Desktop/test.log",os.O_CREATE|os.O_APPEND|os.O_RDWR,0744)
	log.SetOutput(file)
}

func simpleHttpGet(url string) {
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("Error fetching url %s : %s", url, err.Error())
	} else {
		log.Printf("Status Code for %s : %s", url, resp.Status)
		resp.Body.Close()
	}
}