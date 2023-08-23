# 👍 ggogio

[![Go Report Card](https://goreportcard.com/badge/github.com/atgane/ggogio)](https://goreportcard.com/report/github.com/atgane/ggogio) [![golang-test](https://github.com/atgane/ggogio/actions/workflows/test.yml/badge.svg)](https://github.com/atgane/ggogio/actions/workflows/test.yml)

ggogio is a lightweight tcp server framework ... world like to be. It aims for simple use and lightness, maybe.

# 😶‍🌫️ how to use

1. Create a sample chatting client.

```go
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
```

2. Create SampleClient factory.

```go
type SampleFactory struct {
}

func (s SampleFactory) Create() ggogio.Client {
	return new(SampleClient)
}
```

3. Run server!

```go
import (
	"log"

	"github.com/atgane/ggogio"
)

func main() {
	addr := ":10000"
	s := ggogio.NewServer(addr, SampleFactory{})
	err := s.Listen()
	if err != nil {
		log.Fatal(err)
	}
}
```

4. connect with telnet

```sh
telnet localhost 10000
```