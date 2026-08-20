package booking

import (
	"context"
	"errors"
	"time"
)

var (
	ERR_SEAT_ALREADY_TAKEN = errors.New("seat already taken")
)

// Booking represents a confirmed seat reservation.
type Booking struct {
	ID        string
	MovieID   string
	SeatID    string
	UserID    string
	Status    string
	ExpiresAt time.Time
}

type BookingStore interface {
	Book(b Booking) (Booking, error)
	ListBookings(movieID string) []Booking
	RemoveSession(context context.Context, sessionid, userid string) error
	ConfirmSession(context context.Context, sessionid, userid string) (Booking, error)
}
