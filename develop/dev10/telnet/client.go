package telnet

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const bufSize = 4096

type (
	// TelnetClient
	TelnetClient struct {
		Conn    net.Conn
		Timeout time.Duration
		EndCh   chan struct{}
		MsgCh   chan string
		ctx     context.Context
		cancel  context.CancelFunc
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
	// If no connect to the server => break conn on timeout
	dialer := &net.Dialer{Timeout: timeout}
	// conn, err := net.Dial("tcp", fmt.Sprintf("%s:%s", host, port))
	conn, err := dialer.Dial("tcp", fmt.Sprintf("%s:%s", host, port))
	if err != nil {
		return nil, err
	}

	tc := &TelnetClient{
		Conn:    conn,
		Timeout: timeout,
		EndCh:   make(chan struct{}),
		MsgCh:   make(chan string),
	}

	ctx, c := context.WithCancel(context.Background())
	tc.ctx = ctx
	tc.cancel = c
	return tc, nil
}

func (tc *TelnetClient) stop() error {
	tc.cancel()
	time.Sleep(time.Second)
	if err := tc.Close(); err != nil {
		return fmt.Errorf("client connection close error: %w", err)
	}
	return nil
}

func (tc *TelnetClient) WaitOSKill() {
	go func() {
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
		// If receiving ctrl d / ctrl c signal => tc.endchannel will be triggered after that. read/write loop will be linked to endCh
		sig := <-ch
		fmt.Println("\nGot signal:", sig)
		tc.EndCh <- struct{}{}
	}()
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

func (tc *TelnetClient) Cancel() error {
	tc.cancel()
	time.Sleep(time.Second)
	if err := tc.Close(); err != nil {
		return err
	}
	return nil
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

func (tc *TelnetClient) LaunchReadWrite() {
	go tc.readLoop()
	go tc.writeLoop()
}

func (tc *TelnetClient) readLoop() {
	buf := make([]byte, 1024)
	// Infinite loop for handling done signal or write loop
	for {
		select {
		case <-tc.ctx.Done():
			log.Println("read stopped")
			break
		default:
			// Setting deadline for each read attempt
			if err := tc.Conn.SetReadDeadline(time.Now().Add(tc.Timeout)); err != nil {
				log.Println(err)
			}
			// Reading data from socket
			n, err := tc.Conn.Read(buf)
			if err != nil {
				if err == io.EOF {
					log.Println("Remote host aborted connection, exiting from reading...")
					tc.EndCh <- struct{}{}
					break
				}
				if netErr, ok := err.(net.Error); ok && !netErr.Timeout() {
					log.Println(err)
				}
			}
			// If no data is read => no actions
			if n == 0 {
				break
			}
			// If data sent from the socket are ok => print it in stdout
			bs := buf[:n]
			if len(bs) != 0 {
				fmt.Print(string(bs))
			}

		}
	}
}

func (tc *TelnetClient) writeLoop() {
	// Handling input from stdin & sending it to channel
	go func(msgCh chan<- string) {
		reader := bufio.NewReader(os.Stdin)
		for {
			s, err := reader.ReadString('\n')
			if err != nil {
				if err == io.EOF {
					log.Print("Ctrl+D detected, aborting...")
					tc.EndCh <- struct{}{}
					return
				}
				log.Println(err)
			}
			msgCh <- s
		}
	}(tc.MsgCh)

OUTER:
	// If client context is done (this happens when ctrl d detected) => stop writing loop
	for {
		select {
		case <-tc.ctx.Done():
			log.Print("Exiting from writing...")
			break OUTER
		default:
			// Loop for reading from msgChannel (which is storing string from stdin) and writing that data to server
		STDIN:
			for {
				select {
				case msg, ok := <-tc.MsgCh:
					// If message is empty
					if !ok {
						break STDIN
					}
					if _, err := tc.Conn.Write([]byte(msg)); err != nil {
						log.Println(err)
					}
					// Wait deadline for input
				case <-time.After(time.Second):
					break STDIN
				}
			}
		}
	}
	log.Println("...exited from writing")
}
