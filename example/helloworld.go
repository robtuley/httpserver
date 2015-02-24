package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

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
	log.Println("started:> ")

	osStopC := make(chan os.Signal)
	signal.Notify(osStopC,
		syscall.SIGKILL,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
	)

WaitForStop:
	for {
		select {
		case s := <-osStopC:
			log.Println("stopping:> from OS signal", s)
			service.Stop()
		case <-service.HasStoppedC:
			log.Println("stopped:> all conns closed")
			break WaitForStop
		}
	}

	if err = service.Err(); err != nil {
		log.Println("error:> ", err)
	}
}
