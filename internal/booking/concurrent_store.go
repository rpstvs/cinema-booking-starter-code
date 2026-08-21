package booking

import "sync"

// type BookingStore interface {
// 	Book(b Booking) (Booking, error)
// 	ListBookings(movieID string) []Booking
// }

type ConcurrentMemoryStore struct {
	sync.RWMutex
	SeatMap map[string]Booking
}

func NewConcurrentMemoryStore() *ConcurrentMemoryStore {
	return &ConcurrentMemoryStore{
		SeatMap: make(map[string]Booking),
	}
}

func (s *ConcurrentMemoryStore) Book(b Booking) (Booking, error) {
	s.Lock()
	defer s.Unlock()
	_, exists := s.SeatMap[b.SeatID]

	if exists {
		return Booking{}, ERR_SEAT_ALREADY_TAKEN
	}

	s.SeatMap[b.SeatID] = b

	return b, nil
}

func (s *ConcurrentMemoryStore) ListBookings(movieID string) []Booking {
	s.RLock()
	defer s.RUnlock()
	bookingList := make([]Booking, 0, 0)

	for _, v := range s.SeatMap {
		if v.MovieID == movieID {
			bookingList = append(bookingList, v)
		}
	}

	return bookingList
}
