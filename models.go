package main

type Team struct {
	Name         string
	Strength     int
	Points       int
	Wins         int
	Draws        int
	Losses       int
	GoalsFor     int
	GoalsAgainst int
}

type Match struct {
	HomeTeam  string
	AwayTeam  string
	HomeGoals int
	AwayGoals int
	Week      int
}
