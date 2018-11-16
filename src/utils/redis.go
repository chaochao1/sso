package utils

import (
	"github.com/go-redis/redis"
	"sync"
)

var RedisClient *redisClient

type redisClient struct {
	mux   sync.Mutex
	items map[string]*redis.Client
}

type RedisClientInterface interface {
	Set(k string, e *redis.Client)
	Get(k string) (e *redis.Client)
}

func initRedis() {
	RedisClient = NewRedisPool()
	for _, r := range Config.Redis {
		client, _ := r.GetClient()
		RedisClient.Set(r.Name, client)
	}
}

func NewRedisPool() *redisClient {
	r := &redisClient{}
	r.items = make(map[string]*redis.Client)
	return r
}

func (c *redisClient) Set(k string, e *redis.Client) {
	c.mux.Lock()
	c.items[k] = e
	c.mux.Unlock()
}

func (c *redisClient) Get(k string) (e *redis.Client, found bool) {
	c.mux.Lock()
	defer c.mux.Unlock()
	if e, found = c.items[k]; !found {
		return nil, false
	}
	return
}