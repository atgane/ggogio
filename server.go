package ggogio

import (
	"log"
	"net"
	"sync"
)

// Server is a instance to run TCP socket server.
type Server struct {
	Config  interface{}
	factory Factory

	addr         string
	listener     net.Listener
	sessions     []*Session
	sessionsLock sync.RWMutex
	addChan      chan *Session
	removeChan   chan *Session
}

// create new server.
func NewServer(addr string, f Factory) *Server {
	s := new(Server)
	s.factory = f

	s.addr = addr
	s.sessions = []*Session{}
	s.addChan = make(chan *Session, serverAddChanSize)
	s.removeChan = make(chan *Session, serverRemoveChanSize)
	return s
}

// set server config
func (s *Server) SetServerConfig(c interface{}) {
	s.Config = c
}

// get server config
func (s *Server) GetServerConfig() interface{} {
	return s.Config
}

// run TCP socket server.
func (s *Server) Listen() error {
	l, err := net.Listen("tcp", s.addr)
	s.listener = l
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

		session := newSession(done, sendBuf, recvBuf)
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

// find other clients with condition function in specific client.
func (s *Server) Query(condition func(*Session) bool) []*Session {
	s.sessionsLock.RLock()
	defer s.sessionsLock.RUnlock()
	sessions := []*Session{}
	for _, session := range s.sessions {
		if condition(session) {
			sessions = append(sessions, session)
		}
	}
	return sessions
}

func (s *Server) serve() {
	for {
		select {
		case session := <-s.addChan:
			func() {
				s.sessionsLock.Lock()
				defer s.sessionsLock.Unlock()

				s.sessions = append(s.sessions, session)
				log.Printf("client connect: %s\n", session.Addr)
			}()
		case session := <-s.removeChan:
			func() {
				s.sessionsLock.RLock()
				defer s.sessionsLock.RUnlock()
				for i := range s.sessions {
					if s.sessions[i] == session {
						e := len(s.sessions) - 1
						s.sessions[e], s.sessions[i] = s.sessions[i], s.sessions[e]
						s.sessions = s.sessions[:e]
						log.Printf("client leave: %s\n", session.Addr)
						break
					}
				}
			}()
		}
	}
}
