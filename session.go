package ggogio

// Session is a instance to communicate server or other clients.
type Session struct {
	Config interface{}

	Addr string

	recvbuf chan []byte
	sendbuf chan []byte
	done    chan bool
}

func newSession(done chan bool, send chan []byte, recv chan []byte) *Session {
	s := new(Session)
	s.done = done
	s.sendbuf = send
	s.recvbuf = recv
	return s
}

// close session.
func (s *Session) Close() {
	s.done <- true
}

// write data to client.
func (s *Session) Write(request []byte) {
	s.sendbuf <- request
}

// read data from client.
func (s *Session) Read() []byte {
	return <-s.recvbuf
}

// set config
func (s *Session) SetConfig(c interface{}) {
	s.Config = c
}

// get config
func (s *Session) GetConfig() interface{} {
	return s.Config
}
