package main

import (
	"fmt"
	"net/http"
)

func main() {
	InitDB()
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})

	// API endpoints
	http.HandleFunc("/play-week", PlayWeekHandler)
	http.HandleFunc("/play-all", PlayAllHandler)
	http.HandleFunc("/standings", StandingsHandler)
	http.HandleFunc("/results", ResultsHandler)
	http.HandleFunc("/edit-match", EditMatchHandler)
	http.HandleFunc("/predict", PredictHandler)

	fmt.Println("Server running on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
