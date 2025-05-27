package main

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func InitDB() {
	generateSchedule()
	var err error
	DB, err = sql.Open("sqlite3", "./league.db")
	if err != nil {
		panic(err)
	}

	createTables()

	Teams, err = LoadTeams()
	if err != nil {
		panic("Failed to load teams: " + err.Error())
	}

	Matches, err = LoadMatches()
	if err != nil {
		panic("Failed to load matches: " + err.Error())
	}

	for _, m := range Matches {
		if m.Week > CurrentWeek {
			CurrentWeek = m.Week
		}
	}

	if len(Teams) == 0 {
		Teams = []Team{
			{"Liverpool", 90, 0, 0, 0, 0, 0, 0},
			{"Arsenal", 85, 0, 0, 0, 0, 0, 0},
			{"Manchester United", 80, 0, 0, 0, 0, 0, 0},
			{"Chelsea", 75, 0, 0, 0, 0, 0, 0},
		}

		for _, t := range Teams {
			_ = SaveTeam(t)
		}

		fmt.Println("Inserted initial teams into the database.")
	}
}

func createTables() {
	teamTable := `
    CREATE TABLE IF NOT EXISTS teams (
        name TEXT PRIMARY KEY,
        strength INTEGER,
        points INTEGER,
        wins INTEGER,
        draws INTEGER,
        losses INTEGER,
        goals_for INTEGER,
        goals_against INTEGER
    );`

	matchTable := `
    CREATE TABLE IF NOT EXISTS matches (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        home_team TEXT,
        away_team TEXT,
        home_goals INTEGER,
        away_goals INTEGER,
        week INTEGER
    );`

	_, err := DB.Exec(teamTable)
	if err != nil {
		panic(err)
	}

	_, err = DB.Exec(matchTable)
	if err != nil {
		panic(err)
	}
}

func SaveTeam(team Team) error {
	_, err := DB.Exec(`
        INSERT OR REPLACE INTO teams 
        (name, strength, points, wins, draws, losses, goals_for, goals_against)
        VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		team.Name, team.Strength, team.Points, team.Wins, team.Draws, team.Losses, team.GoalsFor, team.GoalsAgainst)
	return err
}

func SaveMatch(match Match) error {
	_, err := DB.Exec(`
        INSERT INTO matches (home_team, away_team, home_goals, away_goals, week)
        VALUES (?, ?, ?, ?, ?)`,
		match.HomeTeam, match.AwayTeam, match.HomeGoals, match.AwayGoals, match.Week)
	return err
}

func LoadTeams() ([]Team, error) {
	rows, err := DB.Query(`SELECT name, strength, points, wins, draws, losses, goals_for, goals_against FROM teams`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var teams []Team
	for rows.Next() {
		var t Team
		rows.Scan(&t.Name, &t.Strength, &t.Points, &t.Wins, &t.Draws, &t.Losses, &t.GoalsFor, &t.GoalsAgainst)
		teams = append(teams, t)
	}
	fmt.Println("Loaded", len(teams), "teams from DB")
	return teams, nil
}

func LoadMatches() ([]Match, error) {
	rows, err := DB.Query(`SELECT home_team, away_team, home_goals, away_goals, week FROM matches`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var matches []Match
	for rows.Next() {
		var m Match
		rows.Scan(&m.HomeTeam, &m.AwayTeam, &m.HomeGoals, &m.AwayGoals, &m.Week)
		matches = append(matches, m)
	}

	return matches, nil
}
