package sshx

import (
	"strings"

	"golang.org/x/crypto/ssh"

	"github.com/xmhrnx/serverus/internal/domain"
)

// keyboardInteractive answers challenges: a lone hidden "Password:"
// prompt gets the password, anything else (one-time codes) is asked.
func (d *Dialer) keyboardInteractive(s domain.Server) ssh.AuthMethod {
	attempt := 0
	return ssh.KeyboardInteractive(func(name, instruction string, questions []string, echos []bool) ([]string, error) {
		answers := make([]string, len(questions))
		for i, q := range questions {
			if len(questions) == 1 && !echos[i] && strings.Contains(strings.ToLower(q), "password") {
				attempt++
				pw, ok := d.opts.Credentials.Password(s, attempt > 1)
				if !ok {
					return nil, errCanceled
				}
				answers[i] = pw
				continue
			}
			prompt := strings.TrimSpace(strings.Join([]string{instruction, q}, " "))
			a, ok := d.opts.Credentials.Answer(s, prompt, echos[i])
			if !ok {
				return nil, errCanceled
			}
			answers[i] = a
		}
		return answers, nil
	})
}
