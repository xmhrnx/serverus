package shell

import (
	"slices"
	"sort"
	"strings"

	"github.com/xmhrnx/serverus/internal/domain"
)

// filtered returns the servers matching the search, grouped and sorted.
func (s *Sidebar) filtered() (groups []string, byGroup map[string][]domain.Server) {
	q := strings.ToLower(strings.TrimSpace(s.search.Text()))
	byGroup = map[string][]domain.Server{}
	for _, srv := range s.env.Repo.Servers() {
		if q != "" && !matches(srv, q) {
			continue
		}
		g := srv.DisplayGroup()
		if _, ok := byGroup[g]; !ok {
			groups = append(groups, g)
		}
		byGroup[g] = append(byGroup[g], srv)
	}
	sort.SliceStable(groups, func(i, j int) bool {
		// "Local" goes last; the rest alphabetically.
		if (groups[i] == "Local") != (groups[j] == "Local") {
			return groups[j] == "Local"
		}
		return strings.ToLower(groups[i]) < strings.ToLower(groups[j])
	})
	for _, g := range groups {
		sort.SliceStable(byGroup[g], func(i, j int) bool {
			a, b := byGroup[g][i], byGroup[g][j]
			if a.Favorite != b.Favorite {
				return a.Favorite
			}
			return strings.ToLower(a.Name) < strings.ToLower(b.Name)
		})
	}
	return groups, byGroup
}

func matches(s domain.Server, q string) bool {
	fields := append([]string{s.Name, s.Host, s.Username, s.Group}, s.Tags...)
	return slices.ContainsFunc(fields, func(f string) bool { return strings.Contains(strings.ToLower(f), q) })
}
