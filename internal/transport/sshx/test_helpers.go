package sshx

import (
	"encoding/binary"

	"golang.org/x/crypto/ssh"
)

func itoa(n uint32) string {
	b := make([]byte, 0, 10)
	if n == 0 {
		return "0"
	}
	for n > 0 {
		b = append([]byte{byte('0' + n%10)}, b...)
		n /= 10
	}
	return string(b)
}

func exitStatus(ch ssh.Channel, code int) {
	var b [4]byte
	binary.BigEndian.PutUint32(b[:], uint32(code))
	ch.SendRequest("exit-status", false, b[:])
}
