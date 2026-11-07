package dialogs

func (h *HostKeyPrompt) answer(ok bool) {
	if h.reply == nil {
		return
	}
	h.open = false
	h.reply <- ok
	h.reply = nil
}
