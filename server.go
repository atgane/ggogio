package ggogio

import (
	"log"
	"net"
)

type Server struct {
	Config  interface{}
	factory Factory

	addr       string
	sessions   []*Session
	addChan    chan *Session
	removeChan chan *Session
}

func NewServer(addr string, f Factory) *Server {
	s := new(Server)
	s.factory = f

	s.addr = addr
	s.sessions = []*Session{}
	s.addChan = make(chan *Session)
	s.removeChan = make(chan *Session)
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
		recvBuf := make(chan []byte, clientDefaultSendChanSize)
		sendBuf := make(chan []byte, clientDefaultRecvChanSize)
		done := make(chan bool, 1)

		session := NewSession(done, sendBuf, recvBuf)
		session.Addr = conn.RemoteAddr().String()

		if err := ic.Init(s, session); err != nil {
			log.Printf("client interface initialize failed: %s\n", err)
		}
		c := newClient(conn, ic, s, session)

		s.addChan <- session

		go c.read()
		go c.write()
	}
}

func (s *Server) serve() {
	for {
		select {
		case session := <-s.addChan:
			s.sessions = append(s.sessions, session)
			log.Printf("client connect: %s\n", session.Addr)
		case session := <-s.removeChan:
			for i := range s.sessions {
				if s.sessions[i] == session {
					s.sessions = append(s.sessions[:i], s.sessions[i+1:]...)
					log.Printf("client leave: %s\n", session.Addr)
					break
				}
			}
		}
	}
}
