package httpserver

import (
	"errors"
	"net"
	"strconv"
	"time"
)

// see http://www.hydrogen18.com/blog/stop-listening-http-server-go.html
type gracefulListener struct {
	*net.TCPListener
	stopC chan bool
}

func newGracefulListener(port int) (*gracefulListener, error) {
	gL := &gracefulListener{
		stopC: make(chan bool),
	}

	listener, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		return gL, err
	}

	tcpL, ok := listener.(*net.TCPListener)
	if !ok {
		return gL, errors.New("Cannot wrap listener")
	}
	gL.TCPListener = tcpL

	return gL, nil
}

func (gL *gracefulListener) Accept() (net.Conn, error) {
	for {
		gL.SetDeadline(time.Now().Add(time.Second))
		newConn, err := gL.TCPListener.Accept()

		select {
		case <-gL.stopC:
			return nil, errors.New("Listener stopped")
		default:
			// still listening, continue as normal
		}

		if err != nil {
			netErr, ok := err.(net.Error)

			// if this is a timeout, then continue to wait for
			// new connections
			if ok && netErr.Timeout() && netErr.Temporary() {
				continue
			}
		}

		return newConn, err
	}
}

func (gL *gracefulListener) Stop() {
	close(gL.stopC)
}
