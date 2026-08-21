package booking

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/sikozonpc/cinema/internal/utils"
)

type holdRequest struct {
	UserID string `json:"userid"`
}

type handler struct {
	svc Service
}

func NewHandler(svc Service) *handler {
	return &handler{
		svc: svc,
	}
}

func (h *handler) ListBookings(w http.ResponseWriter, r *http.Request) {

	movieid := r.PathValue("movieID")
	bookings := h.svc.ListBookings(movieid)

	var seats []seatInfo

	for _, v := range bookings {
		seat := seatInfo{
			SeatInfo: v.SeatID,
			UserID:   v.UserID,
			Booked:   true,
		}
		seats = append(seats, seat)
	}

	utils.WriteJSON(w, http.StatusOK, seats)
}

func (h *handler) HoldSeat(w http.ResponseWriter, r *http.Request) {
	movieid := r.PathValue("movieID")
	seatID := r.PathValue("seatID")

	var req holdRequest
	defer r.Body.Close()
	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, nil)
	}

	data := Booking{
		UserID:  req.UserID,
		MovieID: movieid,
		SeatID:  seatID,
	}
	session, err := h.svc.Book(data)
	if err != nil {
		return
	}

	type holdResponse struct {
		SessionID string `json:"session_id"`
		MovieID   string `json:"movie_id"`
		SeatID    string `json:"seat_id"`
		ExpiresAt string `json:"expires_at"`
	}

	utils.WriteJSON(w, http.StatusOK, holdResponse{
		SessionID: session.ID,
		MovieID:   session.MovieID,
		SeatID:    seatID,
		ExpiresAt: session.ExpiresAt.Format(time.RFC3339),
	})
}

type seatInfo struct {
	SeatInfo string `json:"seat_id"`
	UserID   string `json:"user_id"`
	Booked   bool   `json:"booked"`
}

func (h *handler) ConfirmSession(w http.ResponseWriter, r *http.Request) {

	sessionId := r.PathValue("sessionID")

	var req holdRequest
	defer r.Body.Close()
	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, "bad request")
		return
	}

	if req.UserID == "" {
		utils.WriteJSON(w, http.StatusBadRequest, "userid is mandatory")
		return
	}

	booking, err := h.svc.ConfirmSession(r.Context(), sessionId, req.UserID)

	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, "internal server error")
		return
	}

	confirmation := sessionResponse{
		SessionID: sessionId,
		MovieID:   booking.MovieID,
		SeatID:    booking.SeatID,
		UserID:    req.UserID,
		Status:    booking.Status,
	}

	utils.WriteJSON(w, http.StatusOK, confirmation)
}

type sessionResponse struct {
	SessionID string `json:"session_id"`
	MovieID   string `json:"movie_id"`
	SeatID    string `json:"seat_id"`
	UserID    string `json:"user_id"`
	Status    string `json:"status"`
	ExpiresAt string `json:"expires_at,omitempty"`
}

func (h *handler) RemoveSession(w http.ResponseWriter, r *http.Request) {

	sessionId := r.PathValue("sessionID")

	var req holdRequest
	defer r.Body.Close()
	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, "bad request")
		return
	}

	if req.UserID == "" {
		utils.WriteJSON(w, http.StatusBadRequest, "userid is mandatory")
		return
	}

	err = h.svc.RemoveSession(r.Context(), sessionId, req.UserID)

	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, APIResponse{Error: "internal server error"})
	}

	w.WriteHeader(http.StatusNoContent)
}

type APIResponse struct {
	Data  any    `json:"data"`
	Error string `json:"error,omitempty"`
}
