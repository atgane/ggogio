package ggogio

import (
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type testEchoFactory struct{}

func (t testEchoFactory) Create() Client {
	return new(testEchoClient)
}

var testServerScenario []serverScenario = []serverScenario{
	{"create client and connect test echo server", testServerConnect},
	{"query to find all clients", testServerQuery},
	{"remove client", testServerRemoveClient},
}

type serverScenario struct {
	name string
	fn   func(t *testing.T, state *serverState)
}

type serverState struct {
	server *Server
	client net.Conn
}

type testEchoClient struct {
	server  *Server
	session *Session
}

func (t *testEchoClient) Init(server *Server, session *Session) error {
	t.server = server
	t.session = session

	return nil
}

func (t *testEchoClient) OnLoop() {
	data := t.session.Read()
	t.session.Write(data)
}

func (t *testEchoClient) Close() {
	t.session.Close()
}

func TestServer(t *testing.T) {
	func() {
		state := new(serverState)

		// init test
		addr := ":0"
		s := NewServer(addr, testEchoFactory{})
		state.server = s
		go s.Listen()

		// do test
		for idx := range testServerScenario {
			t.Run(testServerScenario[idx].name, func(t *testing.T) {
				testServerScenario[idx].fn(t, state)
			})
		}
	}()
}

func testServerConnect(t *testing.T, state *serverState) {
	t.Helper()

	time.Sleep(10 * time.Millisecond)
	server := state.server
	serverAddr := server.listener.Addr().String()

	dial, err := net.Dial("tcp", serverAddr)
	require.NoError(t, err)
	state.client = dial

	time.Sleep(10 * time.Millisecond)
	require.Equal(t, 1, len(server.sessions))
	require.Equal(t, 0, len(server.removeChan))
	require.Equal(t, 0, len(server.addChan))
}

func testServerQuery(t *testing.T, state *serverState) {
	t.Helper()

	server := state.server
	sessions := server.Query(func(s *Session) bool { return true })
	require.Equal(t, 1, len(sessions))
}

func testServerRemoveClient(t *testing.T, state *serverState) {
	t.Helper()

	server := state.server
	client := state.client
	err := client.Close()
	require.NoError(t, err)
	time.Sleep(10 * time.Millisecond)
	require.Equal(t, 0, len(server.sessions))
	require.Equal(t, 0, len(server.removeChan))
	require.Equal(t, 0, len(server.addChan))
}
