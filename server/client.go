// This file consist client side of Gache. You can use it in you Go code
// like this:
//      c := server.NewClient(serverAddr)
//      resp, err := c.Sendf("SET key value 10", key, val, ttl)
//      if err != nil {
//          ... 
//      }
//      ...
//
//      resp, err = c.Sendf("GET key value 10", key, val, ttl)
//      if err != nil {
//          ... 
//      }
//      fmt.Printf("Returned value is: %s", resp)

package server

import (
    "fmt"
    "net"
    "bufio"
    "strings"
    // "io/ioutil"
)

// Client is wrapper about net.TCPConn and some validation
type Client struct {
	addr *net.TCPAddr
	Conn *net.TCPConn
	KeepAlive bool
}

// NewClient return pointer to new created Client
func NewClient(addr string) *Client {
	var err error
	c := &Client{KeepAlive: false}
	c.addr, err = net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		panic("Unable resolve addr. Error: " + err.Error())
	}
	return c
}

// Send make new TCP reqiest and return response
func (c *Client) Send(s string) (string, error) {
	defer func() {
		if ! c.KeepAlive {
			c.Close()
		}
	}()

	// Create connection
	var err error
    if ! c.KeepAlive || c.Conn == nil {
		c.Conn, err = net.DialTCP("tcp", nil, c.addr)
		if err != nil {
			log.Warnf("Dial error: %v", err)
			return "", err
		}
    }

	// Write
	_, err = c.Conn.Write([]byte(s + "\r\n"))
	if err != nil {
		log.Debugf("Write error: %s\n", err.Error())
		return "", err
	}
    c.Conn.CloseWrite()

	// Read response
    buf, err := bufio.NewReader(c.Conn).ReadString('\n')
	if err != nil {
		log.Debugf("Read error: %s\n", err.Error())
		return "", err
	}
    c.Conn.CloseRead()
	return strings.TrimSpace(string(buf)), nil
}

// Sendf is shorctut for Send method with parameters subtituting
func (c *Client) Sendf(s string, args ...interface{}) (string, error) {
    msg := fmt.Sprintf(s, args...)
    return c.Send(msg)
}

func (c *Client) Close() {
	c.Conn.Close()
}
