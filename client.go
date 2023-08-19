package ggogio

import (
	"log"
	"net"
)

const (
	clientDefaultSendChanSize = 10
	clientDefaultRecvChanSize = 10
	clientDefaultBufSize      = 1024
)

type client struct {
	Client
	conn     net.Conn
	sendBufs chan []byte
	recvBufs chan []byte
	done     chan bool
}

type Client interface {
	// Init() method is called when Server instance creates Client
	// after tcp connection success
	Init() error

	// Connect() method is ... 아 설명하기 귀찮다...
	Connect(server *Server, recv <-chan []byte, write chan<- []byte) error
	Close() error
}

func newClient(conn net.Conn, ic Client, s *Server) (*client, error) {
	c := new(client)
	c.Client = ic
	c.conn = conn
	c.sendBufs = make(chan []byte, clientDefaultSendChanSize)
	c.recvBufs = make(chan []byte, clientDefaultRecvChanSize)
	c.done = make(chan bool, 1)

	err := c.Client.Connect(s, c.recvBufs, c.sendBufs)
	return c, err
}

func (c *client) read() {
	buf := make([]byte, clientDefaultBufSize)

	for {
		select {
		case <-c.done:
			return
		default:
			n, err := c.conn.Read(buf)
			if err != nil {
				log.Printf("read failed: %s\n", err)
			}
			c.sendBufs <- buf[:n]
		}

		if len(c.sendBufs) == clientDefaultSendChanSize {
			buf := []byte{}
			for i := 0; i < clientDefaultSendChanSize; i++ {
				buf = append(buf, <-c.sendBufs...)
			}
			c.sendBufs <- buf
		}
	}
}

func (c *client) write() {
	for {
		select {
		case <-c.done:
			return
		default:
			buf := <-c.recvBufs

			write := 0
			for write != len(buf) {
				w, err := c.conn.Write(buf)
				if err != nil {
					log.Printf("write failed: %s\n", err)
				}
				write += w
			}
		}
	}
}
