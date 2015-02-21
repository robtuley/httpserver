HTTP Server with Shutdown
=========================

Adds graceful shutdown & health-check utilities to the standard http.Server. 

[![GoDoc](https://godoc.org/github.com/robtuley/httpserver?status.png)](https://godoc.org/github.com/robtuley/httpserver)

Usage
-----

    func main() {
    	service, err := httpserver.New(8080)
    	if err != nil {
    		panic(err)
    	}
    
    	service.HandleFunc("/", func(res http.ResponseWriter, req *http.Request) {
    		res.Header().Set("Content-Type", "text/plain; charset=utf-8")
    		io.WriteString(res, "Hello World")
    	})
    
    	service.Run()
    
    	stopC := make(chan os.Signal)
    	signal.Notify(stopC,
    		syscall.SIGKILL,
    		syscall.SIGQUIT,
    		syscall.SIGHUP,
    		syscall.SIGINT,
    		syscall.SIGTERM,
    	)
    	<-stopC
    
    	service.Stop()
    }
