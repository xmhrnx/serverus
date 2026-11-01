package domain

// Accents are the avatar colors a server can have.
var Accents = []string{
	"#7c6cff", "#22c3a6", "#f2994a", "#eb5757", "#2f80ed",
	"#bb6bd9", "#27ae60", "#f2c94c", "#56ccf2", "#ff7eb6",
}

// AccentFor picks a stable accent for a server without one.
func AccentFor(id string) string {
	var h uint32
	for _, c := range id {
		h = h*31 + uint32(c)
	}
	return Accents[int(h%uint32(len(Accents)))]
}
