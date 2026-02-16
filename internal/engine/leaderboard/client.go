package leaderboard

import (
	"fmt"
	"sort"
)

type Entry struct {
	Username string
	Score    int
}

// GetLeaderboard returns a sorted list of 60 records.
func GetLeaderboard() []Entry {
	entries := []Entry{
		{"mrbhnauka", 65},
		{"Seniru", 300},
		{"Heshan", 200},
		{"Dasun", 150},
		{"Kalindu", 100},
		{"Nathum", 90},
		{"Kalum", 60},
	}

	// Add more dummy records to reach 60
	for i := len(entries) + 1; i <= 60; i++ {
		entries = append(entries, Entry{
			Username: fmt.Sprintf("User %d", i),
			Score:    60 - i, // Ensuring they are sorted by score
		})
	}

	// Double check sorting
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Score > entries[j].Score
	})

	return entries
}
