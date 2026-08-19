package booking

import "context"

type Service struct {
	store BookingStore
}

func NewService(store BookingStore) *Service {
	return &Service{store: store}
}

func (s *Service) Book(b Booking) (Booking, error) {
	book, err := s.store.Book(b)

	if err != nil {
		return Booking{}, err
	}

	return book, nil
}

func (s *Service) ListBookings(movieID string) []Booking {
	bookings := s.store.ListBookings(movieID)

	return bookings
}
func (s *Service) ConfirmSession(context context.Context, sessionId, userid string) (Booking, error) {
	return s.store.ConfirmSession(context, sessionId, userid)

}

func (s *Service) RemoveSession(context context.Context, sessionId, userid string) error {
	return s.store.RemoveSession(context, sessionId, userid)

}
