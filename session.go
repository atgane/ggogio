package ggogio

// 아 설명하기 정말 귀찮다...
type Session struct {
	Config interface{}

	sendBufs chan<- []byte
	recvBufs <-chan []byte
	done     chan<- bool
}

func NewSession(done chan<- bool, recv <-chan []byte, send chan<- []byte) *Session {
	s := new(Session)
	s.done = done
	s.recvBufs = recv
	s.sendBufs = send
	return s
}

func (s *Session) Close() {
	s.done <- true
}

func (s *Session) Write(request []byte) {
	s.sendBufs <- request
}

func (s *Session) Read() []byte {
	return <-s.recvBufs
}

func (s *Session) SetConfig(c interface{}) {
	s.Config = c
}

func (s *Session) GetConfig() interface{} {
	return s.Config
}
