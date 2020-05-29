package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"time"
)

const network = "tcp"

var (
	ErrConnectionClosed = errors.New("connection closed by peer")
	errLog              = log.New(os.Stderr, "", 0)
)

type Client struct {
	address string
	timeout time.Duration
	conn    net.Conn
	input   io.ReadCloser
	output  io.Writer
}

type TelnetClient interface {
	Connect() error
	Receive() error
	Send() error
	Close() error
}

func transferData(in io.Reader, out io.Writer) error {
	sc := bufio.NewScanner(in)
	for sc.Scan() {
		if err := sc.Err(); err != nil {
			return err
		}
		text := sc.Text()
		if text != "" {
			_, err := fmt.Fprintf(out, text+"\n")
			if err != nil {
				errLog.Println("...Connection was closed by peer")
				return ErrConnectionClosed
			}
		}
	}
	return nil
}

func (c *Client) Connect() error {
	conn, err := net.DialTimeout(network, c.address, c.timeout)
	if err != nil {
		return err
	}
	c.conn = conn
	errLog.Println("...Connected to " + c.address)
	return nil
}

func (c *Client) Send() error {
	return transferData(c.input, c.conn)
}

func (c *Client) Receive() error {
	return transferData(c.conn, c.output)
}

func (c *Client) Close() error {
	err := c.input.Close()
	if err != nil {
		return err
	}
	err = c.conn.Close()
	return err
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	return &Client{
		address: address,
		timeout: timeout,
		input:   in,
		output:  out,
	}
}
