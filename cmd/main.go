package main

import (
	"log"
	"net/http"

	go_redis "github.com/sikozonpc/cinema/internal/adapters/redis"
	"github.com/sikozonpc/cinema/internal/booking"
	"github.com/sikozonpc/cinema/internal/utils"
)

const ADDR = ":6379"

func main() {
	store := booking.NewRedisStore(go_redis.NewClient(ADDR))
	svc := booking.NewService(store)
	handler := booking.NewHandler(*svc)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /movies", listMovies)
	mux.HandleFunc("GET /movies/{movieID}/seats", handler.ListBookings)
	mux.HandleFunc("POST /movies/{movieID}/seats/{seatID}/hold", handler.HoldSeat)
	mux.HandleFunc("PUT /sessions/{sessionID}/confirm", handler.ConfirmSession)
	mux.HandleFunc("DELETE /sessions/{sessionID}", handler.RemoveSession)
	mux.Handle("GET /", http.FileServer(http.Dir("static")))
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("failed :%v", err)
	}
}

var moviesResp = []movieResponse{
	{ID: "inception", Title: "Inception", Rows: 5, SeatsPerRow: 8},
	{ID: "dune", Title: "Dune Part Two", Rows: 5, SeatsPerRow: 8},
}

func listMovies(w http.ResponseWriter, r *http.Request) {

	utils.WriteJSON(w, http.StatusOK, moviesResp)
}

type movieResponse struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Rows        int64  `json:"rows"`
	SeatsPerRow int64  `json:"seats_per_row"`
}
