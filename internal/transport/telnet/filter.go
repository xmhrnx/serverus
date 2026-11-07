package telnet

// filter returns the data bytes of p, handling the commands in it. The
// parser state carries over between calls, so commands split across
// reads are handled.
func (n *negotiator) filter(p []byte) []byte {
	out := p[:0]
	var verb byte
	for _, b := range p {
		switch n.state {
		case stData:
			if b == cmdIAC {
				n.state = stIAC
				continue
			}
			out = append(out, b)
		case stIAC:
			switch b {
			case cmdIAC:
				out = append(out, b)
				n.state = stData
			case cmdWill, cmdWont, cmdDo, cmdDont:
				verb = b
				n.sub = append(n.sub[:0], b)
				n.state = stOption
			case cmdSB:
				n.sub = n.sub[:0]
				n.state = stSub
			default:
				n.state = stData
			}
		case stOption:
			verb = n.sub[0]
			n.answer(verb, b)
			n.state = stData
		case stSub:
			if b == cmdIAC {
				n.state = stSubIAC
				continue
			}
			n.sub = append(n.sub, b)
		case stSubIAC:
			if b == cmdSE {
				n.subnegotiation(n.sub)
				n.state = stData
				continue
			}
			n.sub = append(n.sub, b)
			n.state = stSub
		}
	}
	return out
}
