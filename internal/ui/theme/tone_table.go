package theme

// Tones are the background tones to choose from; the first is the
// default, the dark emerald of the app's icon.
var Tones = []Tone{
	{
		ID: "emerald", Name: "Emerald", Terminal: "#09130f", Swatch: "#1b553c",
		Dark:  Surfaces{Background: "#08110d", Surface: "#0c1713", Border: "#1d3028", Text: "#e5eee9", Muted: "#83998e", Secondary: "#97ada2"},
		Light: Surfaces{Background: "#f0f6f3", Surface: "#ffffff", Border: "#d9e6df", Text: "#0f1a15", Muted: "#5f756a", Secondary: "#4f665b"},
	},
	{
		ID: "graphite", Name: "Graphite", Terminal: "#0c0f15", Swatch: "#3a4152",
		Dark:  Surfaces{Background: "#0b0d12", Surface: "#141821", Border: "#232937", Text: "#e8eaf1", Muted: "#8890a5", Secondary: "#9aa3b8"},
		Light: Surfaces{Background: "#f4f5f9", Surface: "#ffffff", Border: "#e2e5ee", Text: "#11141d", Muted: "#69718a", Secondary: "#5b6378"},
	},
	{
		ID: "midnight", Name: "Midnight", Terminal: "#0a1120", Swatch: "#1f3b70",
		Dark:  Surfaces{Background: "#090f1c", Surface: "#111a2c", Border: "#1f2b44", Text: "#e6ebf5", Muted: "#8592ad", Secondary: "#99a6c2"},
		Light: Surfaces{Background: "#f0f4fa", Surface: "#ffffff", Border: "#dbe2ef", Text: "#0f1524", Muted: "#64708a", Secondary: "#55627c"},
	},
	{
		ID: "plum", Name: "Plum", Terminal: "#120d18", Swatch: "#4b2d60",
		Dark:  Surfaces{Background: "#110c16", Surface: "#1a1422", Border: "#2c2238", Text: "#eee8f3", Muted: "#9a8ca8", Secondary: "#ab9dba"},
		Light: Surfaces{Background: "#f6f2f8", Surface: "#ffffff", Border: "#e6deec", Text: "#1a1220", Muted: "#766983", Secondary: "#665a73"},
	},
	{
		ID: "black", Name: "Black", Terminal: "#000000", Swatch: "#000000",
		Dark:  Surfaces{Background: "#000000", Surface: "#0f0f11", Border: "#222226", Text: "#ececee", Muted: "#8b8b94", Secondary: "#a0a0a8"},
		Light: Surfaces{Background: "#f5f5f6", Surface: "#ffffff", Border: "#e4e4e7", Text: "#111113", Muted: "#6b6b74", Secondary: "#5a5a62"},
	},
}

// ToneNamed returns the tone with id, or the default.
func ToneNamed(id string) Tone {
	for _, t := range Tones {
		if t.ID == id {
			return t
		}
	}
	return Tones[0]
}
