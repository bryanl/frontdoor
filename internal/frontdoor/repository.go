package frontdoor

import (
	"context"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

func NewPool(addr string) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			logrus.WithField("redis-addr", addr).Info("redis dial connection")

			c, err := redis.Dial("tcp", addr)
			if err != nil {
				return nil, errors.Wrap(err, "dialing redis connection")
			}

			_, err = c.Do("SELECT", 0)
			if err != nil {
				return nil, errors.Wrap(err, "selecting redis db 0")
			}

			return c, nil
		},
	}
}

type Repository interface {
	AddName(context.Context, string) error
	ListNames(context.Context) ([]string, error)
	Ready(context.Context) error
}

type RedisRepository struct {
	pool   *redis.Pool
	logger logrus.FieldLogger
}

var _ Repository = (*RedisRepository)(nil)

func NewRedisRepository(addr string, logger logrus.FieldLogger) *RedisRepository {
	if addr == "" {
		addr = ":6379"
	}

	rp := NewPool(addr)

	return &RedisRepository{
		logger: logger,
		pool:   rp,
	}
}

func (r *RedisRepository) AddName(ctx context.Context, name string) error {
	conn, err := r.pool.GetContext(ctx)
	if err != nil {
		return errors.Wrap(err, "retrieving redis connection from pool")
	}

	defer conn.Close()

	_, err = conn.Do("RPUSH", "gb", name)
	if err != nil {
		return errors.Wrap(err, "adding name to redis list")
	}

	r.logger.WithFields(logrus.Fields{
		"name": name,
	}).Info("adding name to list")

	return nil
}

func (r *RedisRepository) ListNames(ctx context.Context) ([]string, error) {
	conn, err := r.pool.GetContext(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "retrieving redis connection from pool")
	}

	defer conn.Close()

	names, err := redis.Strings(conn.Do("LRANGE", "gb", 0, -1))
	if err != nil {
		return nil, errors.Wrap(err, "listing guestbook names")
	}

	return names, nil
}

func (r *RedisRepository) Ready(ctx context.Context) error {
	conn, err := r.pool.GetContext(ctx)
	if err != nil {
		return errors.Wrap(err, "retrieving redis connection from pool")
	}

	defer conn.Close()

	resp, err := redis.String(conn.Do("PING"))
	if err != nil {
		return errors.Wrap(err, "pinging redis server")
	}

	if resp != "PONG" {
		return errors.Errorf("expected PONG; server returned %s", resp)
	}

	return nil
}
