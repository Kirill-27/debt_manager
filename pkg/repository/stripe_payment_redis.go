package repository

import (
	"errors"
	"github.com/go-redis/redis"
	"strconv"
	"time"
)

type StripePaymentRedis struct {
	redis *redis.Client
}

const stripeLastHandledTime = "stripeLastHandledTime"

func NewStripePaymentRedis(redis *redis.Client) *StripePaymentRedis {
	return &StripePaymentRedis{redis: redis}
}

func (r *StripePaymentRedis) GetLastHandled() (time.Time, error) {
	stringTime, err := r.redis.Get(stripeLastHandledTime).Result()
	if err != nil {
		if err.Error() == errors.New("redis: nil").Error() {
			return time.Now().AddDate(-1, 0, 0), nil
		}
		return time.Time{}, err
	}
	intTime, err := strconv.ParseInt(stringTime, 10, 64)
	if err != nil {
		return time.Time{}, err
	}
	return time.Unix(intTime, 0), nil
}

func (r *StripePaymentRedis) SetLastHandled(lastHandledTime int64) error {
	return r.redis.Set(stripeLastHandledTime, lastHandledTime, 0).Err()
}
