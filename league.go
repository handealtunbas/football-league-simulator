package main

import (
	"errors"
	"fmt"
	"math/rand"
)

func PlayWeek() {
	if CurrentWeek >= len(scheduledWeeks) {
		fmt.Println("All weeks played.")
		return
	}

	CurrentWeek++
	weekMatches := scheduledWeeks[CurrentWeek-1]

	for _, match := range weekMatches {
		h := findTeamByName(match.HomeTeam)
		a := findTeamByName(match.AwayTeam)

		hGoals := rand.Intn(h.Strength / 20)
		aGoals := rand.Intn(a.Strength / 20)

		if hGoals > aGoals {
			h.Points += 3
			h.Wins++
			a.Losses++
		} else if aGoals > hGoals {
			a.Points += 3
			a.Wins++
			h.Losses++
		} else {
			h.Points++
			a.Points++
			h.Draws++
			a.Draws++
		}

		h.GoalsFor += hGoals
		h.GoalsAgainst += aGoals
		a.GoalsFor += aGoals
		a.GoalsAgainst += hGoals

		played := Match{
			HomeTeam:  h.Name,
			AwayTeam:  a.Name,
			HomeGoals: hGoals,
			AwayGoals: aGoals,
			Week:      CurrentWeek,
		}

		Matches = append(Matches, played)
		SaveTeam(*h)
		SaveTeam(*a)
		SaveMatch(played)
	}
}

func findTeamByName(name string) *Team {
	for i := range Teams {
		if Teams[i].Name == name {
			return &Teams[i]
		}
	}
	return nil
}

func EditMatch(week int, index int, newH int, newA int) error {
	count := 0
	for i, m := range Matches {
		if m.Week == week {
			if count == index {
				// rollback old scores
				for i := range Teams {
					if Teams[i].Name == m.HomeTeam {
						Teams[i].GoalsFor -= m.HomeGoals
						Teams[i].GoalsAgainst -= m.AwayGoals
						Teams[i].Points -= getPoints(m.HomeGoals, m.AwayGoals)
					}
					if Teams[i].Name == m.AwayTeam {
						Teams[i].GoalsFor -= m.AwayGoals
						Teams[i].GoalsAgainst -= m.HomeGoals
						Teams[i].Points -= getPoints(m.AwayGoals, m.HomeGoals)
					}
				}

				// apply new scores
				Matches[i].HomeGoals = newH
				Matches[i].AwayGoals = newA

				for i := range Teams {
					if Teams[i].Name == m.HomeTeam {
						Teams[i].GoalsFor += newH
						Teams[i].GoalsAgainst += newA
						Teams[i].Points += getPoints(newH, newA)
					}
					if Teams[i].Name == m.AwayTeam {
						Teams[i].GoalsFor += newA
						Teams[i].GoalsAgainst += newH
						Teams[i].Points += getPoints(newA, newH)
					}
				}

				return nil
			}
			count++
		}
	}

	return errors.New("match not found")
}

func getPoints(goals1, goals2 int) int {
	if goals1 > goals2 {
		return 3
	} else if goals1 == goals2 {
		return 1
	}
	return 0
}

var scheduledWeeks [][]Match

func generateSchedule() {
	teamCount := len(Teams)
	if teamCount%2 != 0 {
		fmt.Println("Odd number of teams not supported yet")
		return
	}

	var schedule [][]Match
	teamNames := make([]string, teamCount)
	for i, t := range Teams {
		teamNames[i] = t.Name
	}

	fixed := teamNames[0]
	rotating := append([]string{}, teamNames[1:]...)

	rounds := teamCount - 1
	for i := 0; i < rounds; i++ {
		var week []Match

		// Setting home/away for fair fixture
		if i%2 == 0 {
			week = append(week, Match{HomeTeam: fixed, AwayTeam: rotating[0]})
		} else {
			week = append(week, Match{HomeTeam: rotating[0], AwayTeam: fixed})
		}

		for j := 1; j < teamCount/2; j++ {
			home := rotating[j]
			away := rotating[len(rotating)-j]
			if i%2 == 0 {
				week = append(week, Match{HomeTeam: home, AwayTeam: away})
			} else {
				week = append(week, Match{HomeTeam: away, AwayTeam: home})
			}
		}

		schedule = append(schedule, week)
		// Rotate the teams
		rotating = append(rotating[1:], rotating[0])
	}

	// Create reverse fixtures (second leg)
	var secondLeg [][]Match
	for i, week := range schedule {
		var reversedWeek []Match
		for _, match := range week {
			reversedWeek = append(reversedWeek, Match{
				HomeTeam:  match.AwayTeam,
				AwayTeam:  match.HomeTeam,
				HomeGoals: 0,
				AwayGoals: 0,
				Week:      i + rounds + 1, // Set correct week number
			})
		}
		secondLeg = append(secondLeg, reversedWeek)
	}

	scheduledWeeks = append(schedule, secondLeg...)
}
