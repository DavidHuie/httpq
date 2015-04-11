package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/DavidHuie/httpq"
	"github.com/DavidHuie/httpq/Godeps/_workspace/src/github.com/boltdb/bolt"
	"github.com/DavidHuie/httpq/Godeps/_workspace/src/github.com/garyburd/redigo/redis"
	"github.com/DavidHuie/httpq/queue/boltqueue"
	"github.com/DavidHuie/httpq/queue/redisqueue"
)

func main() {
	var dbPath string
	flag.StringVar(&dbPath, "db_path", "httpq.db", "the path to the database file")
	var port string
	flag.StringVar(&port, "port", ":3000", "the port to listen on")
	var redisUrl string
	flag.StringVar(&redisUrl, "redis_url", "", "the url for redis (optional)")
	var redisIdleConnections int
	flag.IntVar(&redisIdleConnections, "redis_idle_connections", 50, "maximum number of idle redis connections (only applicable if using redis)")
	flag.Parse()

	var queue httpq.Queue

	if redisUrl == "" {
		db, err := bolt.Open(dbPath, 0600, nil)
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()
		queue = boltqueue.NewBoltQueue(db)
	} else {
		connManager := NewRedisConnManager(redisUrl)
		redisPool := redis.NewPool(connManager.newRedisConn, redisIdleConnections)
		queue = redisqueue.NewRedisQueue(redisPool)
	}

	hq := httpq.NewHttpq(queue)
	server := httpq.NewServer(hq)

	http.HandleFunc("/push", server.Push)
	http.HandleFunc("/pop", server.Pop)
	http.HandleFunc("/size", server.Size)

	log.Fatal(http.ListenAndServe(port, nil))
}
