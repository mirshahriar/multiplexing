package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"sync"

	mgrpc "github.com/mirshahriar/multiplexing/grpc"
	mhttp "github.com/mirshahriar/multiplexing/http"
)

type bufferData struct {
	// As we check the first 16 bytes of the connection to identify if it's gRPC or not,
	// we need to store the data in a buffer to be able to read it again from gRPC/http server
	buffer []byte
	// bufferRead is the number of bytes read from the buffer
	bufferRead int
	bufferUsed bool
}

func (s *muxConn) ReadBufferLine() func() (string, error) {
	var line string
	var err error
	var once sync.Once
	return func() (string, error) {
		once.Do(func() {
			b := make([]byte, 16)
			var n int
			n, err = s.Conn.Read(b)
			if err == nil {
				line = strings.Split(string(b), "\n")[0]
			}

			s.buf.buffer = make([]byte, n)
			s.buf.bufferRead = 0
			copy(s.buf.buffer, b[:n])
		})
		return line, err
	}
}

func (s *muxConn) Read(p []byte) (int, error) {
	// If we have remaining data in the buffer, return it first
	if len(s.buf.buffer) > s.buf.bufferRead && !s.buf.bufferUsed {
		bn := copy(p, s.buf.buffer[s.buf.bufferRead:])
		s.buf.bufferRead += bn
		return bn, nil
	}

	sn, sErr := s.Conn.Read(p)
	// when read from buffer is done, we don't need the buffer anymore
	s.buf.bufferUsed = true
	s.buf.buffer = nil

	return sn, sErr
}

type muxConn struct {
	net.Conn
	buf     *bufferData
	bufLine func() (string, error)
}

func newMuxConn(c net.Conn) *muxConn {
	mc := &muxConn{
		Conn: c,
		buf:  &bufferData{},
	}
	mc.bufLine = mc.ReadBufferLine()

	return mc
}

type muxListener struct {
	net.Listener
	conn chan *muxConn
}

func (l *muxListener) Accept() (net.Conn, error) {
	return <-l.conn, nil
}

func main() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		_ = listener.Close()
	}()

	grpcConn := &muxListener{
		Listener: listener,
		conn:     make(chan *muxConn, 1024),
	}
	httpConn := &muxListener{
		Listener: listener,
		conn:     make(chan *muxConn, 1024),
	}

	grpcServer := mgrpc.NewGRPCServer()
	httpServer := mhttp.NewHTTPServer()

	defer func() {
		grpcServer.GracefulStop()
		_ = httpServer.Shutdown(nil)
	}()

	go func() {
		fmt.Println("gRPC Server is running on port 8080")
		if err := grpcServer.Serve(grpcConn); err != nil {
			log.Fatal(err)
		}
		fmt.Println("gRPC Server is closed")
		os.Exit(1)
	}()

	go func() {
		fmt.Println("Server is running on port 8080")
		if err := httpServer.Serve(httpConn); err != nil {
			log.Fatal(err)
		}
		fmt.Println("Server is closed")
		os.Exit(1)
	}()

	for {
		c, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}

		mc := newMuxConn(c)

		headLine, err := mc.bufLine()
		if err != nil {
			log.Fatal(err)
		}

		isGrpc := isGRPC(headLine)

		if isGrpc {
			grpcConn.conn <- mc
		} else {
			httpConn.conn <- mc
		}

		fmt.Println("New connection is accepted")
	}

}

func isGRPC(headLine string) bool {
	return strings.HasPrefix(headLine, "PRI * HTTP/2.0")
}
