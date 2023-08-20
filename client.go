package ggogio

import (
	"io"
	"log"
	"net"
)

type client struct {
	Client
	conn    net.Conn
	recvBuf chan []byte
	sendBuf chan []byte
	done    chan bool
}

type Client interface {
	// Init() method is called when Server instance creates Client
	// after tcp connection success
	Init(*Server, *Session) error

	// OnLoop() method is ... 아 설명하기 귀찮다...
	OnLoop()

	// Close() method is called when Client called Session.Close().
	// implement termination connection using this function.
	Close()
}

func newClient(conn net.Conn, ic Client, s *Server, session *Session) *client {
	c := new(client)
	c.Client = ic
	c.conn = conn
	c.recvBuf = session.recvbuf
	c.sendBuf = session.sendbuf
	c.done = session.done

	go c.onLoop()

	return c
}

func (c *client) close() {
	for len(c.sendBuf) > 0 {
	}
	c.conn.Close()
	c.Client.Close()
}

func (c *client) onLoop() {
	for {
		select {
		case <-c.done:
			c.close()
			return
		default:
			c.Client.OnLoop()
		}
	}
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
				if err == io.EOF {
					log.Printf("client connection closed: %s\n", err)
					c.close()
					return
				}
				log.Printf("read failed: %s\n", err)
			}
			c.recvBuf <- buf[:n]
		}

		if len(c.recvBuf) == clientDefaultSendChanSize {
			buf := []byte{}
			for i := 0; i < clientDefaultSendChanSize; i++ {
				buf = append(buf, <-c.recvBuf...)
			}
			c.recvBuf <- buf
		}
	}
}

func (c *client) write() {
	for {
		select {
		case <-c.done:
			return
		default:
			buf := <-c.sendBuf

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
