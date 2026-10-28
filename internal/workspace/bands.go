package workspace

import (
	"errors"
	"fmt"
)

// BandSuggestions are the bands the CLI offers a person to confirm or replace:
// one for each member who holds none, and one for the next member added.
type BandSuggestions struct {
	Members map[string]int `json:"members,omitempty"`
	Next    int            `json:"next"`
}

// SuggestBands offers bands that start empty. It never offers band 0, never a
// band another member holds or one already offered in the same answer, and
// never one whose intent or plan range already holds a taken number, so the
// member given it allocates from the first number of its range. A person may
// still pick any of those by hand; allocation skips the taken numbers.
func (s *Store) SuggestBands() (BandSuggestions, error) {
	suggestions := BandSuggestions{Members: map[string]int{}}
	members, err := s.Members()
	if err != nil {
		return suggestions, err
	}
	unavailable := map[int]bool{}
	for _, member := range members.Members {
		if member.Band != nil {
			unavailable[*member.Band] = true
		}
	}
	for _, kind := range []string{"intent", "plan"} {
		taken, err := s.takenNumbers(kind)
		if err != nil {
			return suggestions, fmt.Errorf("bands cannot be suggested while the taken numbers cannot be read: %w", err)
		}
		for n := range taken {
			unavailable[n/bandWidth(kind)] = true
		}
	}
	candidate := 0
	next := func() (int, error) {
		for candidate++; candidate <= 9999; candidate++ {
			if !unavailable[candidate] {
				return candidate, nil
			}
		}
		return 0, errors.New("every band from 1 to 9999 is held or already holds records")
	}
	for _, id := range members.IDs() {
		if members.Members[id].Band != nil {
			continue
		}
		band, err := next()
		if err != nil {
			return suggestions, err
		}
		suggestions.Members[id] = band
	}
	if suggestions.Next, err = next(); err != nil {
		return suggestions, err
	}
	return suggestions, nil
}
