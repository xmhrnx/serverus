package settings

import (
	"slices"

	"github.com/xmhrnx/giode"
	"github.com/xmhrnx/giode/components/dropdown"

	"github.com/xmhrnx/serverus/internal/domain"
)

func selectOf(values []int, current int, label func(int) string) *dropdown.Select {
	opts := make([]string, len(values))
	for i, v := range values {
		opts[i] = label(v)
	}
	s := giode.Input.Select(selectStyles, opts...)
	s.SetSelected(max(0, slices.Index(values, current)))
	return s
}

func (p *Page) update(fn func(*domain.Settings)) {
	s := p.env.Repo.Settings()
	fn(&s)
	if err := p.env.Repo.SaveSettings(s); err != nil {
		p.env.Fail("Save settings", err)
		return
	}
	if p.onChange != nil {
		p.onChange(p.env.Repo.Settings())
	}
}
