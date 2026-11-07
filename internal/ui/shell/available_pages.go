package shell

import (
	"github.com/xmhrnx/serverus/internal/domain"
)

// availablePages lists the pages a protocol supports.
func availablePages(p domain.Protocol) []pageSpec {
	caps := p.Capabilities()
	var out []pageSpec
	for _, s := range pageSpecs {
		if caps.Has(s.needs) && (s.localOK || p != domain.ProtocolLocal) {
			out = append(out, s)
		}
	}
	return out
}
