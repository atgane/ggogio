package main

import (
	"fmt"

	"github.com/atgane/ggogio"
)

type SampleClient struct {
	server  *ggogio.Server
	session *ggogio.Session
}

func (s *SampleClient) Init(server *ggogio.Server, session *ggogio.Session) error {
	s.server = server
	s.session = session

	sessions := s.server.Query(func(session *ggogio.Session) bool {
		return s.session != session
	})
	for _, session := range sessions {
		session.Write([]byte(fmt.Sprintf("%s joined. \n", s.session.Addr)))
	}

	return nil
}

func (s *SampleClient) OnLoop() {
	data := s.session.Read()
	sessions := s.server.Query(func(session *ggogio.Session) bool {
		return s.session != session
	})
	for _, session := range sessions {
		session.Write(data)
	}
}

func (s *SampleClient) Close() {
	sessions := s.server.Query(func(session *ggogio.Session) bool {
		return s.session != session
	})
	for _, session := range sessions {
		session.Write([]byte(fmt.Sprintf("%s left. \n", s.session.Addr)))
	}
}
