package cache

import "github.com/garyburd/redigo/redis"

type Cache struct {
	url    string
	prefix string
}

func NewCache(ip string, port string, prefix string) *Cache {
	cache := new(Cache)
	cache.url = ip + ":" + port
	cache.prefix = prefix
	return cache
}

func (c *Cache) open() (redis.Conn, error) {
	redisC, err := redis.Dial("tcp", c.url)
	return redisC, err
}

func (c *Cache) SETNX(tag string, key string, value string, t int) error {
	redisC, err := c.open()
	defer redisC.Close()
	if err != nil {
		return err
	}

	skey := c.prefix + tag + key
	_, err = redisC.Do("SETNX", skey, value)
	if err == nil {
		redisC.Do("EXPIRE", skey, t)
	}
	return err
}

func (c *Cache) GET(tag string, key string) (string, error) {
	redisC, err := c.open()
	defer redisC.Close()
	if err != nil {
		return "", err
	}

	return redis.String(redisC.Do("GET", c.prefix+tag+key))
}

func (c *Cache) ClearAll() error {
	redisC, err := c.open()
	defer redisC.Close()
	if err != nil {
		return err
	}

	_, err = redisC.Do("FLUSHALL")
	return err
}
