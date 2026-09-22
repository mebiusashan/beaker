package cache

import (
	"time"

	"github.com/garyburd/redigo/redis"
)

type Cache struct {
	prefix string
	pool   *redis.Pool
}

func NewCache(ip string, port string, prefix string) *Cache {
	cache := &Cache{prefix: prefix}
	cache.pool = &redis.Pool{
		MaxIdle:     10,
		MaxActive:   100,
		IdleTimeout: 5 * time.Minute,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", ip+":"+port,
				redis.DialConnectTimeout(3*time.Second),
				redis.DialReadTimeout(3*time.Second),
				redis.DialWriteTimeout(3*time.Second),
			)
		},
		TestOnBorrow: func(conn redis.Conn, lastUsed time.Time) error {
			if time.Since(lastUsed) < time.Minute {
				return nil
			}
			_, err := conn.Do("PING")
			return err
		},
	}
	return cache
}

func (c *Cache) SETNX(tag string, key string, value string, t int) error {
	redisC := c.pool.Get()
	defer redisC.Close()
	if err := redisC.Err(); err != nil {
		return err
	}
	skey := c.prefix + tag + key
	_, err := redisC.Do("SET", skey, value, "EX", t, "NX")
	return err
}

func (c *Cache) GET(tag string, key string) (string, error) {
	redisC := c.pool.Get()
	defer redisC.Close()
	if err := redisC.Err(); err != nil {
		return "", err
	}
	return redis.String(redisC.Do("GET", c.prefix+tag+key))
}

// ClearAll removes only Beaker's namespaced keys. It never flushes a shared
// Redis instance used by another application.
func (c *Cache) ClearAll() error {
	redisC := c.pool.Get()
	defer redisC.Close()
	if err := redisC.Err(); err != nil {
		return err
	}
	var cursor int64
	pattern := c.prefix + "*"
	for {
		values, err := redis.Values(redisC.Do("SCAN", cursor, "MATCH", pattern, "COUNT", 200))
		if err != nil {
			return err
		}
		var keys []string
		if _, err := redis.Scan(values, &cursor, &keys); err != nil {
			return err
		}
		if len(keys) > 0 {
			args := make([]interface{}, len(keys))
			for i := range keys {
				args[i] = keys[i]
			}
			if _, err := redisC.Do("DEL", args...); err != nil {
				return err
			}
		}
		if cursor == 0 {
			return nil
		}
	}
}
