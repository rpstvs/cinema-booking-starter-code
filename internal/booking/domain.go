package booking

import "errors"

var (
	ERR_SEAT_ALREADY_TAKEN = errors.New("seat already taken")
)

// Booking represents a confirmed seat reservation.
type Booking struct {
	ID      string
	MovieID string
	SeatID  string
	UserID  string
	Status  string
}

type BookingStore interface {
	Book(b Booking) (Booking, error)
	ListBookings(movieID string) []Booking
}
