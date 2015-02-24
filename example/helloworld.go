package main

import (
	"io"
	"log"
	"net/http"

	"github.com/robtuley/httpserver"
)

func main() {
	service, err := httpserver.New(8080)
	if err != nil {
		log.Println("error:> ", err)
		return
	}

	service.HandleFunc("/", func(res http.ResponseWriter, req *http.Request) {
		io.WriteString(res, "Hello World")
	})

	service.Start()
	log.Println("started:> ok")

	sig := service.WaitStop()
	log.Println("stopped:> ", sig)

	if err = service.Err(); err != nil {
		log.Println("error:> ", err)
	}
}
