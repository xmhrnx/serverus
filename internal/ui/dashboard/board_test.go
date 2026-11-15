package dashboard

import (
	"testing"

	"github.com/xmhrnx/giode"

	"github.com/xmhrnx/serverus/internal/domain"
	"github.com/xmhrnx/serverus/internal/storage"
	"github.com/xmhrnx/serverus/internal/ui/env"
)

func kinds(cs []domain.DashboardCard) []string {
	out := make([]string, len(cs))
	for i, c := range cs {
		out[i] = c.Kind + ":" + string(c.Size)
	}
	return out
}

func TestBoardLayoutIsSaved(t *testing.T) {
	repo, err := storage.OpenRepository(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	e := &env.Env{Repo: repo}
	card := func(kind string, size domain.CardSize) Spec {
		return Spec{Kind: kind, Title: kind, Size: size, Render: func() giode.Element { return giode.Box(giode.Styles{}) }}
	}
	specs := []Spec{card("a", domain.CardSmall), card("b", domain.CardMedium), card("c", domain.CardLarge)}
	b := New(e, "test", specs, "a", "b", "gone")

	if got := kinds(b.Cards()); len(got) != 2 || got[0] != "a:small" || got[1] != "b:medium" {
		t.Fatalf("defaults %v", got)
	}
	b.addCard("c")
	b.move("c", -1)
	b.resize("a", domain.CardLarge)
	b.remove("b")
	want := []string{"a:large", "c:large"}
	got := kinds(New(e, "test", specs, "a", "b").Cards()) // a fresh board reads what was saved
	if len(got) != len(want) || got[0] != want[0] || got[1] != want[1] {
		t.Fatalf("saved %v, want %v", got, want)
	}
	b.save(nil)
	if got := kinds(b.Cards()); len(got) != 2 || got[0] != "a:small" {
		t.Fatalf("after reset %v", got)
	}
}

func TestSpan(t *testing.T) {
	for _, c := range []struct {
		size  domain.CardSize
		width int
		want  int
	}{
		{domain.CardSmall, 1200, 3}, {domain.CardMedium, 1200, 6}, {domain.CardLarge, 1200, 12},
		{domain.CardSmall, 700, 6}, {domain.CardMedium, 700, 12},
		{domain.CardSmall, 380, 12},
	} {
		if got := span(c.size, c.width); got != c.want {
			t.Errorf("%s at %d: %d, want %d", c.size, c.width, got, c.want)
		}
	}
}
