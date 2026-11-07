package theme

// Highlights are the highlight colors to choose from; the first is the
// default.
var Highlights = []Highlight{
	{ID: "violet", Name: "Violet", Dark: "#7c6cff", Light: "#6352f5"},
	{ID: "emerald", Name: "Emerald", Dark: "#10b981", Light: "#059669"},
	{ID: "teal", Name: "Teal", Dark: "#14b8a6", Light: "#0d9488"},
	{ID: "blue", Name: "Blue", Dark: "#3b82f6", Light: "#2563eb"},
	{ID: "sky", Name: "Sky", Dark: "#38bdf8", Light: "#0284c7"},
	{ID: "pink", Name: "Pink", Dark: "#ec4899", Light: "#db2777"},
	{ID: "rose", Name: "Rose", Dark: "#f43f5e", Light: "#e11d48"},
	{ID: "orange", Name: "Orange", Dark: "#f97316", Light: "#ea580c"},
	{ID: "amber", Name: "Amber", Dark: "#f59e0b", Light: "#d97706"},
}

// HighlightNamed returns the accent with id, or the default.
func HighlightNamed(id string) Highlight {
	for _, a := range Highlights {
		if a.ID == id {
			return a
		}
	}
	return Highlights[0]
}
