package ggogio

import (
	"io"
	"log"
	"net"
)

type client struct {
	Client
	server  *Server
	session *Session
	conn    net.Conn
	recvBuf chan []byte
	sendBuf chan []byte
	done    chan bool
}

// Client is an interface for handling the connection
// of the client connected to the server.
type Client interface {
	// Init() method is called when Server instance creates Client
	// after tcp connection success
	Init(*Server, *Session) error

	// OnLoop() method is called repeatedly and asynchronously
	// when Init method end.
	OnLoop()

	// Close() method is called when Client called Session.Close().
	// implement termination connection using this function.
	Close()
}

func newClient(conn net.Conn, ic Client, server *Server, session *Session) *client {
	c := new(client)
	c.Client = ic
	c.server = server
	c.session = session
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
	c.server.removeChan <- c.session
	c.Client.Close()
	log.Print("client connection closed\n")
}

func (c *client) onLoop() {
	for {
		select {
		case <-c.done:
			c.done <- true
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
			c.done <- true
			return
		default:
			n, err := c.conn.Read(buf)
			if err != nil {
				log.Printf("read failed. close client connection: %s\n", err)
				c.close()
				return
			}
			c.recvBuf <- buf[:n]
		}

		// when c.recvBuf maxed
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
			c.done <- true
			return
		default:
			buf := <-c.sendBuf

			write := 0
			for write != len(buf) {
				w, err := c.conn.Write(buf)
				if err != nil {
					if err == io.EOF {
						c.close()
						return
					}
					log.Printf("write failed: %s\n", err)
				}
				write += w
			}
		}
	}
}
