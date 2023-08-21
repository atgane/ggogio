package ggogio

import (
	"testing"
)

var testClientScenario []clientScenario = []clientScenario{}

type clientScenario struct {
	name string
	fn   func(t *testing.T, state *clientState)
}

type clientState struct {
	client *client
}

func TestClient(t *testing.T) {
	func() {
		state := new(clientState)

		// init test

		// do test
		for idx := range testClientScenario {
			t.Run(testClientScenario[idx].name, func(t *testing.T) {
				testClientScenario[idx].fn(t, state)
			})
		}
	}()
}

type testClient struct {
	server  *Server
	session *Session
}

func (t *testClient) Init(server *Server, session *Session) error {
	t.server = server
	t.session = session

	return nil
}

func (t *testClient) OnLoop() {
	data := t.session.Read()
	t.session.Write(data)
}

func (t *testClient) Close() {
	t.session.Close()
}
