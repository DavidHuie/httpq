package redisqueue

import "github.com/DavidHuie/httpq/Godeps/_workspace/src/github.com/garyburd/redigo/redis"

const (
	redisListName = "httpq"
)

type RedisQueue struct {
	pool *redis.Pool
}

func NewRedisQueue(pool *redis.Pool) *RedisQueue {
	return &RedisQueue{pool}
}

func (r *RedisQueue) Push(bytes []byte) error {
	conn := r.pool.Get()
	defer conn.Close()

	_, err := conn.Do("LPUSH", redisListName, bytes)
	if err != nil {
		return err
	}

	return nil
}

func (r *RedisQueue) Pop() ([]byte, error) {
	conn := r.pool.Get()
	defer conn.Close()

	response, err := conn.Do("RPOP", redisListName)
	if err != nil {
		return nil, err
	}
	bytes, err := redis.Bytes(response, nil)
	if err != redis.ErrNil && err != nil {
		return nil, err
	}

	return bytes, nil
}

func (r *RedisQueue) Size() (uint64, error) {
	conn := r.pool.Get()
	defer conn.Close()

	response, err := conn.Do("LLEN", redisListName)
	if err != nil {
		return 0, err
	}
	size, err := redis.Uint64(response, nil)
	if err != redis.ErrNil && err != nil {
		return 0, err
	}

	return size, nil
}
