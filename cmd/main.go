package main

import (
	"log"
	"net/http"

	"github.com/sikozonpc/cinema/internal/utils"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /movies", listMovies)
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("failed :v", err)
	}
}

var moviesResp = []movieResponse{
	{ID: "inception", Title: "Inception", Rows: 5, SeatsPerRow: 8},
	{ID: "dune", Title: "Dune Part Two", Rows: 5, SeatsPerRow: 8},
}

func listMovies(w http.ResponseWriter, r *http.Request) {

	utils.WriteJSON(w, http.StatusOK, moviesResp)
}
