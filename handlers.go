package main

import (
	"fmt"
	"net/http"
)
import "strconv"

func PlayWeekHandler(w http.ResponseWriter, r *http.Request) {
	PlayWeek()
	w.Write([]byte("Week played!"))
}

func StandingsHandler(w http.ResponseWriter, r *http.Request) {
	sorted := SortTeams(Teams)
	fmt.Fprintln(w, "Team\tPTS\tP\tW\tD\tL\tGF\tGA\tGD")

	for _, t := range sorted {
		played := t.Wins + t.Draws + t.Losses
		gd := t.GoalsFor - t.GoalsAgainst
		fmt.Fprintf(w, "%s\t%d\t%d\t%d\t%d\t%d\t%d\t%d\t%d\n",
			t.Name, t.Points, played, t.Wins, t.Draws, t.Losses, t.GoalsFor, t.GoalsAgainst, gd)
	}

}

func ResultsHandler(w http.ResponseWriter, r *http.Request) {
	weekStr := r.URL.Query().Get("week")
	if weekStr == "" {
		http.Error(w, "Missing week parameter", http.StatusBadRequest)
		return
	}

	week, err := strconv.Atoi(weekStr)
	if err != nil {
		http.Error(w, "Invalid week", http.StatusBadRequest)
		return
	}

	var result []string
	for _, match := range Matches {
		if match.Week == week {
			line := fmt.Sprintf("%s %d - %d %s", match.HomeTeam, match.HomeGoals, match.AwayGoals, match.AwayTeam)
			result = append(result, line)
		}
	}

	if len(result) == 0 {
		fmt.Fprintf(w, "No matches found for week %d", week)
	} else {
		for _, r := range result {
			fmt.Fprintln(w, r)
		}
	}
}

func PlayAllHandler(w http.ResponseWriter, r *http.Request) {
	for CurrentWeek < len(scheduledWeeks) {
		PlayWeek()
	}
	w.Write([]byte("All weeks played!"))
}

func EditMatchHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	weekStr := r.FormValue("week")
	indexStr := r.FormValue("index")
	newHomeGoalsStr := r.FormValue("homeGoals")
	newAwayGoalsStr := r.FormValue("awayGoals")

	week, _ := strconv.Atoi(weekStr)
	index, _ := strconv.Atoi(indexStr)
	newHomeGoals, _ := strconv.Atoi(newHomeGoalsStr)
	newAwayGoals, _ := strconv.Atoi(newAwayGoalsStr)

	err := EditMatch(week, index, newHomeGoals, newAwayGoals)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Write([]byte("Match updated"))
}

func PredictHandler(w http.ResponseWriter, r *http.Request) {
	if CurrentWeek < 4 {
		fmt.Fprintln(w, "Championship Win Probabilities (available after Week 4)")
		return
	}

	sorted := SortTeams(Teams)
	total := 0
	for _, t := range sorted {
		total += t.Points
	}

	if total == 0 {
		fmt.Fprintln(w, "No games played yet. Cannot generate predictions.")
		return
	}

	fmt.Fprintln(w, "Championship Win Probabilities (after Week", CurrentWeek, "):")
	runningTotal := 0.0
	var percentages []float64

	// Calculate raw percentages
	for _, t := range sorted {
		p := (float64(t.Points) / float64(total)) * 100
		percentages = append(percentages, p)
		runningTotal += p
	}

	// Normalize to exactly 100%
	correction := 100.0 - runningTotal
	percentages[len(percentages)-1] += correction

	for i, t := range sorted {
		fmt.Fprintf(w, "%s – %.0f%%\n", t.Name, percentages[i])
	}
}
