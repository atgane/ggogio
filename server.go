package ggogio

import (
	"log"
	"net"
)

type Server struct {
	Config  interface{}
	factory Factory

	addr       string
	clients    []*client
	addChan    chan *client
	removeChan chan *client
}

func NewServer(addr string, f Factory) *Server {
	s := new(Server)
	s.factory = f

	s.addr = addr
	s.clients = []*client{}
	s.addChan = make(chan *client)
	s.removeChan = make(chan *client)
	return s
}

func (s *Server) SetServerConfig(c interface{}) {
	s.Config = c
}

func (s *Server) GetServerConfig() interface{} {
	return s.Config
}

func (s *Server) Listen() error {
	l, err := net.Listen("tcp", s.addr)
	if err != nil {
		return err
	}

	defer l.Close()

	go s.serve()

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Printf("connection error: %s\n", err)
		}

		ic := s.factory.Create()

		if err := ic.Init(s); err != nil {
			log.Printf("client interface initialize failed: %s\n", err)
		}
		c := newClient(conn, ic, s)

		s.addChan <- c

		go c.read()
		go c.write()
	}
}

func (s *Server) serve() {
	for {
		select {
		case c := <-s.addChan:
			s.clients = append(s.clients, c)
			log.Printf("client connect: %s\n", c.conn.RemoteAddr().String())
		case c := <-s.removeChan:
			for i := range s.clients {
				if s.clients[i] == c {
					s.clients = append(s.clients[:i], s.clients[i+1:]...)
					log.Printf("client leave: %s\n", c.conn.RemoteAddr().String())
					break
				}
			}
		}
	}
}
