package main

import (
	"fmt"
	"log"
	"net"
	"time"
)

const bufSize = 4096

type (
	// TelnetClient
	TelnetClient struct {
		conn    net.Conn
		timeout time.Duration
	}
)

// CTOR for telnet client.
func NewTelnetClient(host, port string, timeout time.Duration) (*TelnetClient, error) {

	// conn, err := net.Dial("tcp", fmt.Sprintf("%s:%s", host, port))
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%s", host, port))
	if err != nil {
		return nil, err
	}

	return &TelnetClient{
		conn:    conn,
		timeout: 10 * time.Second,
	}, nil
}

// Close
func (tc *TelnetClient) Close() error {
	return tc.conn.Close()
}

// Write
func (tc *TelnetClient) Write(in []byte) error {

	// Config request execution parameters.
	if err := tc.conn.SetDeadline(time.Now().Add(tc.timeout)); err != nil {
		return fmt.Errorf("set conn deadline error: %w", err)
	}

	// Write request to connection.
	if _, err := tc.conn.Write(in); err != nil {
		return fmt.Errorf("conn write error: %w", err)
	}

	return nil
}

func (tc *TelnetClient) Read() ([]byte, error) {
	// Read response from connection.
	buf := make([]byte, bufSize)

	if err := tc.conn.SetDeadline(time.Now().Add(tc.timeout)); err != nil {
		return nil, fmt.Errorf("set conn read deadline error: %w", err)
	}

	n, err := tc.conn.Read(buf)
	if err != nil {
		/*
			if v, ok := err.(net.Error); ok && v.Timeout() {
				return nil, fmt.Errorf("read timeout")
			}

			if errors.Is(err, io.EOF) {
				return nil, fmt.Errorf("connection is closed by server")
			}
		*/
		return nil, fmt.Errorf("conn read error: %v", err)
	}

	return buf[:n], nil
}

func main() {

	tc, err := NewTelnetClient("google.com", "telnets", time.Second*5)

	if err != nil {
		log.Fatal(err)
	}

	err = tc.Write([]byte("\n"))

	if err != nil {
		log.Fatal(err)
	}

	resp, err := tc.Read()

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(resp))

	err = tc.Write([]byte("qqck"))

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(resp))

}
