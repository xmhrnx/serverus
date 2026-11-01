// Package telnet opens plain-text terminal sessions to devices that
// only speak Telnet: switches, routers, embedded boards.
package telnet

// Telnet command and option codes (RFC 854, 857, 858, 1073, 1091).
const (
	cmdSE   = 240
	cmdSB   = 250
	cmdWill = 251
	cmdWont = 252
	cmdDo   = 253
	cmdDont = 254
	cmdIAC  = 255

	optEcho     = 1
	optSGA      = 3 // suppress go-ahead
	optTermType = 24
	optNAWS     = 31 // window size

	termTypeIs   = 0
	termTypeSend = 1
)
