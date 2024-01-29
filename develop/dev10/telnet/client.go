package telnet

import (
	"errors"
	"fmt"
	"io"
	"net"
	"time"
)

const bufSize = 4096

type (
	// TelnetClient
	TelnetClient struct {
		Conn    net.Conn
		Timeout time.Duration
	}
)

func Dial(address string, port string) (*TelnetClient, error) {

	const network = "tcp"
	fullAddr := fmt.Sprintf("%s:%s", address, port)
	connection, err := net.Dial(network, fullAddr)
	if err != nil {
		return nil, fmt.Errorf("Could not set dial: %w", err)
	}
	/*
		clientConn := telnetConn{
			conn: connection,
		}
	*/
	client := TelnetClient{Conn: connection, Timeout: time.Second * 5}
	return &client, nil
}

// CTOR for telnet client.
func NewTelnetClient(host, port string, timeout time.Duration) (*TelnetClient, error) {

	// conn, err := net.Dial("tcp", fmt.Sprintf("%s:%s", host, port))
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%s", host, port))
	if err != nil {
		return nil, err
	}

	return &TelnetClient{
		Conn:    conn,
		Timeout: timeout,
	}, nil
}

// Close
func (tc *TelnetClient) Close() error {
	return tc.Conn.Close()
}

// Write
func (tc *TelnetClient) Write(in []byte) error {

	// Config request execution parameters.
	if err := tc.Conn.SetWriteDeadline(time.Now().Add(tc.Timeout)); err != nil {
		return fmt.Errorf("set conn deadline error: %w", err)
	}

	// Write request to connection.
	if _, err := tc.Conn.Write(in); err != nil {
		return fmt.Errorf("conn write error: %w", err)
	}

	return nil
}

func (tc *TelnetClient) SetTimeout(timeout int) {
	if timeout > 0 {
		tc.Timeout = time.Second * time.Duration(timeout)
	}
}

func (tc *TelnetClient) Read() ([]byte, error) {
	// Read response from connection.
	buf := make([]byte, bufSize)
	/*
		if err := tc.Conn.SetDeadline(time.Now().Add(tc.Timeout)); err != nil {
			return nil, fmt.Errorf("set conn read deadline error: %w", err)
		}
	*/
	n, err := tc.Conn.Read(buf)
	if err != nil {
		fmt.Println(err)
		if v, ok := err.(net.Error); ok && v.Timeout() {
			return nil, fmt.Errorf("read timeout")
		}

		if errors.Is(err, io.EOF) {
			return nil, fmt.Errorf("connection is closed by server")
		}

		return nil, fmt.Errorf("conn read error: %v", err)
	}

	return buf[:n], nil
}

/*
func main() {

	tc, err := NewTelnetClient("yahoo.com", "8085", time.Second*5)

	if err != nil {
		log.Fatal(err)
	}
	defer tc.Close()

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
*/
