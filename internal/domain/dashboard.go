package domain

// CardSize is how much of a dashboard row a card takes.
type CardSize string

const (
	CardSmall  CardSize = "small"  // a quarter of a wide row
	CardMedium CardSize = "medium" // half
	CardLarge  CardSize = "large"  // the whole row
)

// Span is the card's width in twelfths of a row.
func (s CardSize) Span() int {
	switch s {
	case CardSmall:
		return 3
	case CardLarge:
		return 12
	}
	return 6
}

// DashboardCard is one card placed on a dashboard.
type DashboardCard struct {
	Kind string   `json:"kind"`
	Size CardSize `json:"size"`
}

// Dashboard names the dashboards whose cards the user arranges.
const (
	DashboardFleet  = "fleet"  // every server at once
	DashboardServer = "server" // one server's overview, shared by all servers
)
