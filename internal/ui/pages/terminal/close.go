package terminal

// Close ends every shell, when the server is removed.
func (p *Page) Close() {
	for _, t := range p.tabs {
		t.close()
	}
	p.tabs = nil
}
