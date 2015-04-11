package boltqueue

import (
	"encoding/binary"
	"encoding/json"

	"github.com/boltdb/bolt"
)

type BoltQueue struct {
	conn *bolt.DB
}

func NewBoltQueue(conn *bolt.DB) *BoltQueue {
	return &BoltQueue{conn}
}

type queueMetadata struct {
	Head uint64
	Last uint64
}

var (
	metadataBucketName = []byte("b")
	metadataKey        = []byte("m")
	dataBucketName     = []byte("d")
)

func getMetadata(b *bolt.Bucket) (*queueMetadata, error) {
	var metadata queueMetadata
	value := b.Get(metadataKey)

	// Create metadata if it doesn't exist
	if value == nil {
		metadata = queueMetadata{}
		bytes, err := json.Marshal(metadata)
		if err != nil {
			return nil, err
		}
		if err := b.Put(metadataKey, bytes); err != nil {
			return nil, err
		}
	} else {
		if err := json.Unmarshal(value, &metadata); err != nil {
			return nil, err
		}
	}

	return &metadata, nil
}

func (b *BoltQueue) Push(bytes []byte) error {
	return b.conn.Update(func(tx *bolt.Tx) error {
		mbucket, err := tx.CreateBucketIfNotExists(metadataBucketName)
		if err != nil {
			return err
		}
		dbucket, err := tx.CreateBucketIfNotExists(dataBucketName)
		if err != nil {
			return err
		}

		metadata, err := getMetadata(mbucket)
		if err != nil {
			return err
		}

		// Update metadata to reflect new data location
		metadata.Last += 1

		// Update recently initialized metadatas
		if metadata.Head == 0 {
			metadata.Head = 1
		}

		metadataBytes, err := json.Marshal(metadata)
		if err != nil {
			return err
		}
		if err := mbucket.Put(metadataKey, metadataBytes); err != nil {
			return err
		}

		// Push
		dataLocationBytes := make([]byte, 8)
		binary.PutUvarint(dataLocationBytes, metadata.Last)
		if err := dbucket.Put(dataLocationBytes, bytes); err != nil {
			return err
		}

		return nil
	})
}

func (b *BoltQueue) Pop() ([]byte, error) {
	var response []byte
	err := b.conn.Update(func(tx *bolt.Tx) error {
		mbucket, err := tx.CreateBucketIfNotExists(metadataBucketName)
		if err != nil {
			return err
		}
		dbucket, err := tx.CreateBucketIfNotExists(dataBucketName)
		if err != nil {
			return err
		}

		metadata, err := getMetadata(mbucket)
		if err != nil {
			return err
		}

		// Handle an empty queue
		if metadata.Head > metadata.Last {
			response = nil
			return nil
		}

		// Perform pop
		dataLocationBytes := make([]byte, 8)
		binary.PutUvarint(dataLocationBytes, metadata.Head)
		response = dbucket.Get(dataLocationBytes)
		if response == nil {
			return nil
		}
		if err := dbucket.Delete(dataLocationBytes); err != nil {
			return err
		}

		// Update metadata
		metadata.Head = metadata.Head + 1
		metadataBytes, err := json.Marshal(metadata)
		if err != nil {
			return err
		}
		if err := mbucket.Put(metadataKey, metadataBytes); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (b *BoltQueue) Size() (uint64, error) {
	var size uint64
	err := b.conn.View(func(tx *bolt.Tx) error {
		mbucket, err := tx.CreateBucketIfNotExists(metadataBucketName)
		if err != nil {
			return err
		}
		metadata, err := getMetadata(mbucket)
		if err != nil {
			return err
		}

		if metadata.Head > metadata.Last {
			size = 0
			return nil
		}

		size = metadata.Last - metadata.Head
		return nil
	})
	if err != nil {
		return 0, err
	}

	return size, nil
}
