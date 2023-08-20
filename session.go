package ggogio

// 아 설명하기 정말 귀찮다...
type Session struct {
	Config interface{}

	Addr string

	recvbuf chan []byte
	sendbuf chan []byte
	done    chan bool
}

func NewSession(done chan bool, send chan []byte, recv chan []byte) *Session {
	s := new(Session)
	s.done = done
	s.sendbuf = send
	s.recvbuf = recv
	return s
}

func (s *Session) Close() {
	s.done <- true
}

func (s *Session) Write(request []byte) {
	s.sendbuf <- request
}

func (s *Session) Read() []byte {
	return <-s.recvbuf
}

func (s *Session) SetConfig(c interface{}) {
	s.Config = c
}

func (s *Session) GetConfig() interface{} {
	return s.Config
}
