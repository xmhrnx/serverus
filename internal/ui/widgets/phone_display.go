package widgets

import "github.com/mehran9675/giode"

// HiddenOnPhones is a Display value that removes an element on phones
// (the compact layout) and keeps it elsewhere.
func HiddenOnPhones() giode.Display {
	if giode.Compact() {
		return giode.DisplayNone
	}
	return ""
}

// ShownOnPhones is a Display value that keeps an element only on phones.
func ShownOnPhones() giode.Display {
	if giode.Compact() {
		return ""
	}
	return giode.DisplayNone
}
