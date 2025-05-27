package main

import "sort"

func SortTeams(teams []Team) []Team {
	sorted := make([]Team, len(teams))
	copy(sorted, teams)

	sort.Slice(sorted, func(i, j int) bool {
		if sorted[i].Points != sorted[j].Points {
			return sorted[i].Points > sorted[j].Points
		}
		goalDiffI := sorted[i].GoalsFor - sorted[i].GoalsAgainst
		goalDiffJ := sorted[j].GoalsFor - sorted[j].GoalsAgainst
		if goalDiffI != goalDiffJ {
			return goalDiffI > goalDiffJ
		}
		return sorted[i].GoalsFor > sorted[j].GoalsFor
	})

	return sorted
}
