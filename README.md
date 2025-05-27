# ⚽ Football League Simulator

A web-based simulator for a fictional football league, implemented in Go with an HTML/JavaScript frontend. The system models team strengths, simulates weekly matches, calculates standings, and displays win probabilities.

##  Features

- Simulates a full league between four teams
- Tracks:
  - Points
  - Wins / Draws / Losses
  - Goals For / Against / Difference
- Weekly match results
- Predicts championship probabilities
- Persistent data using SQLite
- Responsive frontend interface with interactive buttons

## Technologies

- **Backend**: Go (Golang)
- **Frontend**: HTML + Vanilla JavaScript + CSS
- **Database**: SQLite3

## How It Works

1. Each team plays every other team both home and away.
2. Results are randomly generated based on team strength.
3. The table is updated each week.
4. Predictions are generated after Week 4 using Monte Carlo simulations.

## Project Structure

.
├── main.go # Entry point for the web server
├── db.go # SQLite database setup and queries
├── models.go # Data models (Team, Match)
├── league.go # Core simulation logic
├── handlers.go # HTTP request handlers
├── utils.go # Utility functions
├── data.go # Initial data setup
├── index.html # Frontend interface
├── league.db # SQLite database (auto-generated)

## API Endpoints

Method	Endpoint	Description
GET	/standings	Returns the current league standings
GET	/results?week=n	Returns match results for week n
POST	/play-week	Simulates matches for the next week
POST	/play-all	Simulates all remaining matches
GET	/predict	Returns probability of each team winning the league
