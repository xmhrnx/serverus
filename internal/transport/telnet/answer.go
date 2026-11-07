package telnet

func (n *negotiator) answer(verb, opt byte) {
	switch verb {
	case cmdDo:
		switch opt {
		case optTermType, optSGA:
			n.send(cmdIAC, cmdWill, opt)
		case optNAWS:
			n.send(cmdIAC, cmdWill, opt)
			n.mu.Lock()
			n.naws = true
			n.mu.Unlock()
			n.sendSize()
		default:
			n.send(cmdIAC, cmdWont, opt)
		}
	case cmdWill:
		switch opt {
		case optEcho, optSGA:
			n.send(cmdIAC, cmdDo, opt)
		default:
			n.send(cmdIAC, cmdDont, opt)
		}
	}
}
