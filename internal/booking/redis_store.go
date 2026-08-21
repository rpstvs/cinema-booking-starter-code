package booking

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

// type BookingStore interface {
// 	Book(b Booking) (Booking, error)
// 	ListBookings(movieID string) []Booking
// }

const defaultHoldTTL = 2 * time.Minute

type RedisStore struct {
	rdb *redis.Client
}

func NewRedisStore(rdb *redis.Client) *RedisStore {
	return &RedisStore{
		rdb: rdb,
	}
}

func (s *RedisStore) Book(b Booking) (Booking, error) {

	book, err := s.hold(b)

	if err != nil {
		return Booking{}, nil
	}

	return book, nil
}

func (s *RedisStore) ListBookings(movieID string) []Booking {
	pattern := fmt.Sprintf("session:%s:*", movieID)

	var Sessions []Booking

	ctx := context.Background()

	iter := s.rdb.Scan(ctx, 0, pattern, 0).Iterator()

	for iter.Next(ctx) {
		val, err := s.rdb.Get(ctx, iter.Val()).Result()

		if err != nil {
			continue
		}

		session, err := parseSession(val)

		if err != nil {
			continue
		}
		Sessions = append(Sessions, session)
	}
	return Sessions
}

func (s *RedisStore) hold(b Booking) (Booking, error) {
	id := uuid.New().String()
	now := time.Now()
	key := fmt.Sprintf("seat:%s:%s", b.MovieID, b.SeatID)
	b.ID = id
	val, _ := json.Marshal(b)

	ctx := context.Background()

	res := s.rdb.SetArgs(ctx, key, val, redis.SetArgs{
		Mode: "NX",
		TTL:  defaultHoldTTL,
	})

	if res.Val() != "OK" {
		return Booking{}, ERR_SEAT_ALREADY_TAKEN
	}

	s.rdb.Set(ctx, SessionKey(id), key, defaultHoldTTL)

	return Booking{
		ID:        id,
		MovieID:   b.MovieID,
		SeatID:    b.SeatID,
		UserID:    b.UserID,
		Status:    "held",
		ExpiresAt: now.Add(defaultHoldTTL),
	}, nil
}

func (s *RedisStore) ConfirmSession(ctx context.Context, sessionid, userid string) (Booking, error) {
	session, sk, err := s.getSession(ctx, sessionid, userid)

	if err != nil {
		return Booking{}, err
	}

	s.rdb.Persist(ctx, sk)
	s.rdb.Persist(ctx, SessionKey(sessionid))

	session.Status = "confirmed"

	data := Booking{
		ID:      sessionid,
		SeatID:  session.SeatID,
		MovieID: session.MovieID,
		UserID:  session.UserID,
		Status:  session.Status,
	}

	return data, nil
}

func (s *RedisStore) RemoveSession(ctx context.Context, sessionid, userid string) error {
	_, sk, err := s.getSession(ctx, sessionid, userid)

	if err != nil {
		return err
	}

	s.rdb.Del(ctx, sk, SessionKey(sessionid))

	return nil
}

func (s *RedisStore) getSession(ctx context.Context, sessionId, userid string) (Booking, string, error) {
	sk, err := s.rdb.Get(ctx, SessionKey(sessionId)).Result()

	if err != nil {
		return Booking{}, "", err
	}

	val, err := s.rdb.Get(ctx, sk).Result()

	if err != nil {
		return Booking{}, "", err
	}

	booking, err := parseSession(val)

	if err != nil {
		return Booking{}, "", err
	}

	return booking, sk, nil
}

func SessionKey(id string) string {
	return fmt.Sprintf("session:%s", id)
}

func parseSession(val string) (Booking, error) {
	var data Booking

	err := json.Unmarshal([]byte(val), &data)

	if err != nil {
		return Booking{}, err
	}

	return Booking{
		ID:        data.ID,
		MovieID:   data.MovieID,
		SeatID:    data.SeatID,
		UserID:    data.UserID,
		Status:    data.Status,
		ExpiresAt: data.ExpiresAt,
	}, nil
}
