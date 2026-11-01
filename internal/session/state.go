// Package session keeps one live connection per server: connecting and
// reconnecting, metrics collection, and access to the connection's
// capabilities.
package session

// State is a session's connection state.
type State int

const (
	Disconnected State = iota
	Connecting
	Connected
	Failed
)

// Label is the state's display name.
func (s State) Label() string {
	switch s {
	case Connecting:
		return "Connecting"
	case Connected:
		return "Connected"
	case Failed:
		return "Connection failed"
	}
	return "Disconnected"
}
