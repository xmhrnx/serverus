package snippets

import (
	"context"
	"strconv"
	"time"

	"github.com/xmhrnx/serverus/internal/domain"
	"github.com/xmhrnx/serverus/internal/ui/env"
)

func (p *Page) runHere(s domain.Snippet) {
	sess := p.session
	p.output.Show("play", s.Name, p.session.Server().Name+"  ·  "+s.Command, func(ctx context.Context) (string, error) {
		ctx, cancel := context.WithTimeout(ctx, 5*time.Minute)
		defer cancel()
		out, err := runOn(ctx, sess, s.Command)
		if err != nil {
			return "", err
		}
		text := combined(out)
		if out.ExitCode != 0 {
			text += "\n\n[exit status " + strconv.Itoa(out.ExitCode) + "]"
		}
		return text, nil
	})
}

func (p *Page) remove(s domain.Snippet) {
	p.env.Confirm(env.Confirmation{
		Title: "Delete “" + s.Name + "”?", Message: "The snippet is removed from every server's list.",
		Action: "Delete", Destructive: true, OnConfirm: func() { _ = p.env.Repo.DeleteSnippet(s.ID) },
	})
}
