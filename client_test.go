package ggogio

import (
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

var testClientScenario []clientScenario = []clientScenario{
	{"create client and connect test echo server", testClientConnect},
	{"get client session using query to server", testClientQuery},
	{"read hello when client send hello", testClientReadHelloData},
	{"write hello from server to client", testClientWriteHelloData},
	{"close client session", testClientClose},
}

type clientScenario struct {
	name string
	fn   func(t *testing.T, state *clientState)
}

type clientState struct {
	server  *Server
	session *Session
	client  net.Conn
}

func TestClient(t *testing.T) {
	func() {
		state := new(clientState)

		// init test
		addr := ":0"
		s := NewServer(addr, testNoneFactory{})
		state.server = s
		go s.Listen()

		// do test
		for idx := range testClientScenario {
			t.Run(testClientScenario[idx].name, func(t *testing.T) {
				testClientScenario[idx].fn(t, state)
			})
		}
	}()
}

type testNoneFactory struct{}

func (t testNoneFactory) Create() Client {
	return new(testNoneClient)
}

type testNoneClient struct {
	server  *Server
	session *Session
}

func (t *testNoneClient) Init(server *Server, session *Session) error {
	t.server = server
	t.session = session

	return nil
}

func (t *testNoneClient) OnLoop() {
}

func (t *testNoneClient) Close() {
	t.session.Close()
}

func testClientConnect(t *testing.T, state *clientState) {
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

func testClientQuery(t *testing.T, state *clientState) {
	t.Helper()

	server := state.server
	sessions := server.Query(func(s *Session) bool { return true })
	require.Equal(t, 1, len(sessions))
	state.session = sessions[0]
}

func testClientReadHelloData(t *testing.T, state *clientState) {
	t.Helper()

	client := state.client
	n, err := client.Write(hello)
	require.NoError(t, err)
	require.Equal(t, len(hello), n)

	time.Sleep(10 * time.Millisecond)
	session := state.session
	require.Equal(t, 1, len(session.recvbuf))
	read := session.Read()
	require.Equal(t, hello, read)
}

func testClientWriteHelloData(t *testing.T, state *clientState) {
	t.Helper()

	client := state.client
	session := state.session
	session.Write(hello)

	buf := make([]byte, len(hello))
	n, err := client.Read(buf)
	require.NoError(t, err)
	require.Equal(t, len(hello), n)
	require.Equal(t, hello, buf[:n])
}

func testClientClose(t *testing.T, state *clientState) {
	t.Helper()

	server := state.server
	session := state.session
	session.Close()
	require.Equal(t, 1, len(session.done))

	time.Sleep(10 * time.Millisecond)
	sessions := server.Query(func(s *Session) bool { return true })
	require.Equal(t, 0, len(sessions))
}
