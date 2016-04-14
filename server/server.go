package server

import (
	"fmt"
	"net"
	"bufio"
)

// Server is main struct consist method to manage TCP income connections
// and writing response.
type Server struct {
    conn *net.Conn
	addr string
	KeepAlive bool
}

// NewServer create and return Server instance
func NewServer(addr string) *Server {
	return &Server{addr: addr}
}

// ListenTCP starts listen TCP connections
func (s *Server) ListenTCP() {
	ln, err := net.Listen("tcp", s.addr)
	if err != nil {
		log.Fatalf("listen error: %v", err)
	}
	log.Debugf("Server started at %s", s.addr)
    conns := s.clientConns(ln)
	for {
        go s.handleConn(<-conns)
	}
}

func (s *Server) Stop() {
    (*s.conn).Close()
}

// clientConns starts waiting client connections and write them
// to the channed for next handler
func (s *Server) clientConns(ln net.Listener) chan net.Conn {
    ch := make(chan net.Conn)
    go func() {
        for {
            conn, err := ln.Accept()
            if conn == nil {
                log.Errorf("Couldn't accept: " + err.Error())
                continue
            }
            ch <- conn
        }
    }()
    return ch
}

// handleConn get one client connection read, validate and write response
func (s *Server) handleConn(conn net.Conn) {
	defer func() {
		if s.KeepAlive { return }
		conn.Close()
	}()
    b := bufio.NewReader(conn)
    for {
        line, err := b.ReadBytes('\n')
        if err != nil {
            log.Warnf("Error reading bytes: %s", err.Error())
			writeErr(conn, 500, err)
            if s.KeepAlive { continue }
			break
        }

		r, err := NewRequest(string(line))
		if err != nil {
			writeErr(conn, 400, err)
			if s.KeepAlive { continue }
			break
		}

		resp := r.Route()
		if resp.Error != nil {
			writeErr(conn, 400, resp.Error)
			if s.KeepAlive { continue }
			break
		}
		conn.Write([]byte(resp.Body + "\n"))
        if !s.KeepAlive { break }
    }
}

// writeErr is shotrcut for response error writing
func writeErr(conn net.Conn, code int, err error) {
	conn.Write([]byte(fmt.Sprintf("[%d] %s\n", code, err.Error())))
}
