package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/DavidHuie/httpq"
	"github.com/DavidHuie/httpq/Godeps/_workspace/src/github.com/boltdb/bolt"
	"github.com/DavidHuie/httpq/queue/boltqueue"
)

func main() {
	var dbPath string
	flag.StringVar(&dbPath, "db_path", "httpq.db", "the path to the database file")
	var port string
	flag.StringVar(&port, "port", ":3000", "the port to listen on")
	flag.Parse()

	db, err := bolt.Open(dbPath, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	queue := boltqueue.NewBoltQueue(db)
	hq := httpq.NewHttpq(queue)
	server := httpq.NewServer(hq)

	http.HandleFunc("/push", server.Push)
	http.HandleFunc("/pop", server.Pop)
	http.HandleFunc("/size", server.Size)

	log.Fatal(http.ListenAndServe(port, nil))
}
