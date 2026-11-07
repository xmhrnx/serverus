package dialogs

func (p *SecretPrompt) answer(ok bool) {
	if !p.open || p.reply == nil {
		return
	}
	p.open = false
	p.reply <- SecretAnswer{Value: p.field.Text(), Remember: p.req.Remember && p.remember.Checked(), OK: ok}
	p.reply = nil
	p.field.SetText("")
}
