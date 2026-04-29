package booking

// type BookingStore interface {
// 	Book(b Booking) (Booking, error)
// 	ListBookings(movieID string) []Booking
// }

type MemoryStore struct {
	SeatMap map[string]Booking
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		SeatMap: make(map[string]Booking),
	}
}

func (s *MemoryStore) Book(b Booking) (Booking, error) {

	_, exists := s.SeatMap[b.SeatID]

	if exists {
		return Booking{}, ERR_SEAT_ALREADY_TAKEN
	}

	s.SeatMap[b.SeatID] = b

	return b, nil
}

func (s *MemoryStore) ListBookings(movieID string) []Booking {

	bookingList := make([]Booking, 0, 0)

	for _, v := range s.SeatMap {
		if v.MovieID == movieID {
			bookingList = append(bookingList, v)
		}
	}

	return bookingList
}
