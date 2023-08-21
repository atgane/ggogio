package ggogio

import (
	"testing"

	"github.com/stretchr/testify/require"
)

var hello = []byte("hello")
var testConf = testSessionConfig{
	v: 1,
	s: "hello",
}

var testSessionScenario []sessionScenario = []sessionScenario{
	{"write some data to session", testSessionWrite},
	{"read some data from session", testSessionRead},
	{"set config at session", testSessionSetConfig},
	{"get config from session", testSessionGetConfig},
	{"close session", testSessionClose},
}

type sessionScenario struct {
	name string
	fn   func(t *testing.T, state *sessionState)
}

type sessionState struct {
	session *Session
}

type testSessionConfig struct {
	v int
	s string
}

func TestSession(t *testing.T) {
	func() {
		state := new(sessionState)

		// init test
		done := make(chan bool, 1)
		send := make(chan []byte, 10)
		recv := make(chan []byte, 10)
		session := newSession(done, send, recv)
		state.session = session

		// do test
		for idx := range testSessionScenario {
			t.Run(testSessionScenario[idx].name, func(t *testing.T) {
				testSessionScenario[idx].fn(t, state)
			})
		}
	}()
}

func testSessionWrite(t *testing.T, state *sessionState) {
	t.Helper()

	session := state.session
	session.Write(hello)
	require.Equal(t, hello, <-session.sendbuf)
}

func testSessionRead(t *testing.T, state *sessionState) {
	t.Helper()

	session := state.session
	session.recvbuf <- hello
	data := session.Read()
	require.Equal(t, hello, data)
	require.Equal(t, 0, len(session.recvbuf))
}

func testSessionSetConfig(t *testing.T, state *sessionState) {
	t.Helper()

	session := state.session
	session.SetConfig(testConf)
	require.NotNil(t, session.Config)
}

func testSessionGetConfig(t *testing.T, state *sessionState) {
	t.Helper()

	session := state.session
	resultConf := session.GetConfig()
	require.Equal(t, testConf, resultConf.(testSessionConfig))
}

func testSessionClose(t *testing.T, state *sessionState) {
	t.Helper()

	session := state.session
	session.Close()
	require.Equal(t, 1, len(session.done))
	require.Equal(t, true, <-session.done)
}
