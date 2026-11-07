package dialogs

import (
	"github.com/xmhrnx/serverus/internal/domain"
)

// Answer asks a keyboard-interactive question, such as a one-time code.
func (c *Credentials) Answer(s domain.Server, prompt string, echo bool) (string, bool) {
	a := c.Prompt.Ask(SecretRequest{
		Title:   s.Name + " asks",
		Message: prompt,
		Label:   "Answer",
		Hidden:  !echo,
	})
	return a.Value, a.OK
}
