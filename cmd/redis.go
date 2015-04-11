package main

import "github.com/DavidHuie/httpq/Godeps/_workspace/src/github.com/garyburd/redigo/redis"

type redisConnManager struct {
	url string
}

func NewRedisConnManager(url string) *redisConnManager {
	return &redisConnManager{url}
}

func (r *redisConnManager) newRedisConn() (redis.Conn, error) {
	c, err := redis.Dial("tcp", r.url)
	if err != nil {
		return nil, err
	}
	return c, nil
}
