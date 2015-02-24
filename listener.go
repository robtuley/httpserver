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
	HasStoppedC chan bool
	StoppedErr  error
}

func newGracefulListener(port int) (*gracefulListener, error) {
	portStr := ":" + strconv.Itoa(port)
	gL := &gracefulListener{
		HasStoppedC: make(chan bool),
		StoppedErr:  errors.New(portStr + " listener stopped"),
	}

	listener, err := net.Listen("tcp", portStr)
	if err != nil {
		return gL, err
	}

	tcpL, ok := listener.(*net.TCPListener)
	if !ok {
		return gL, errors.New("cannot wrap listener")
	}
	gL.TCPListener = tcpL

	return gL, nil
}

func (gL *gracefulListener) Accept() (net.Conn, error) {
	for {
		gL.SetDeadline(time.Now().Add(time.Second))
		newConn, err := gL.TCPListener.Accept()

		select {
		case <-gL.HasStoppedC:
			return nil, gL.StoppedErr
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
	close(gL.HasStoppedC)
}
