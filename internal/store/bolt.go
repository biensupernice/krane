package store

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/boltdb/bolt"
	"github.com/sirupsen/logrus"

	"github.com/biensupernice/krane/internal/collection"
)

type BoltDB struct {
	*bolt.DB
}

type BoltConfig struct {
	Path    string
	Buckets []string
}

var (
	once     sync.Once
	instance *BoltDB

	defaultBoltPath = "/tmp/krane.db"

	// 0600 - Sets permissions so that:
	// (U)ser / owner can read, can write and can't execute.
	// (G)roup can't read, can't write and can't execute.
	// (O)thers can't read, can't write and can't execute.
	fileMode os.FileMode = 0600
)

func Instance() Store { return instance }

func New(path string) *BoltDB {
	if instance != nil {
		logrus.Info("Bolt instance already exists...")
		return instance
	}

	logrus.Info("Opening store...")

	options := &bolt.Options{Timeout: 30 * time.Second}

	if path == "" {
		path = defaultBoltPath
	}

	db, err := bolt.Open(path, fileMode, options)

	if err != nil {
		logrus.Fatalf("Failed to open to store at %s, %s", path, err.Error())
		return nil
	}

	once.Do(func() { instance = &BoltDB{db} })

	logrus.Infof("Successfully opened store at %s", path)

	// Buckets to create
	bkts := []string{
		collection.Authentication,
		collection.Deployments,
		collection.Sessions,
	}
	instance.createBkts(bkts)

	return instance
}

func (b *BoltDB) createBkts(bkts []string) {
	logrus.Infof("Creating %d bucket(s)", len(bkts))
	for _, bkt := range bkts {
		err := b.Update(func(tx *bolt.Tx) error {
			logrus.Infof("Creating %s bucket", bkt)
			_, err := tx.CreateBucket([]byte(bkt))
			if err != nil {
				return err
			}
			logrus.Infof("Successfully created %s bucket", bkt)
			return nil
		})

		if err != nil {
			logrus.Infof("Unable to create %s bucket: %s", bkt, err)
		}
	}
}

func (b *BoltDB) Shutdown() { b.Close() }

func (b *BoltDB) Put(collection string, key string, value []byte) error {
	return instance.Update(func(tx *bolt.Tx) error {
		bkt := tx.Bucket([]byte(collection))
		if bkt == nil {
			return errors.New(fmt.Sprintf("Bucket %s does not exists", collection))
		}

		return bkt.Put([]byte(key), value)
	})
}

func (b *BoltDB) Get(collection, key string) (data []byte, err error) {
	instance.View(func(tx *bolt.Tx) error {
		bkt := tx.Bucket([]byte(collection))
		if bkt == nil {
			return errors.New(fmt.Sprintf("Bucket %s does not exists", collection))
		}

		data = bkt.Get([]byte(key))
		return nil
	})

	if err != nil {
		return
	}

	return
}

func (b *BoltDB) GetAll(collection string) (data [][]byte, err error) {
	err = instance.View(func(tx *bolt.Tx) (err error) {
		bkt := tx.Bucket([]byte(collection))
		if bkt == nil {
			return errors.New(fmt.Sprintf("Bucket %s does not exists", collection))
		}

		bkt.ForEach(func(k, v []byte) (err error) {
			data = append(data, v)
			return
		})

		return
	})

	if err != nil {
		return
	}

	return
}

// Iterate over a time range.
// minDate: RFC3339 sortable time string ie. 1990-01-01T00:00:00Z
// maxDate example: RFC3339 sortable time string ie. 2000-01-01T00:00:00Z
func (b *BoltDB) GetInRange(collection, minDate, maxDate string) (data [][]byte, err error) {
	err = instance.View(func(tx *bolt.Tx) (err error) {
		bkt := tx.Bucket([]byte(collection))
		if bkt == nil {
			return nil
		}

		c := bkt.Cursor()

		for k, v := c.Seek([]byte(minDate)); k != nil && bytes.Compare(k, []byte(maxDate)) <= 0; k, v = c.Next() {
			data = append(data, v)
		}

		return
	})

	return
}

func (b *BoltDB) Remove(collection string, key string) error {
	return instance.Update(func(tx *bolt.Tx) error {
		bkt := tx.Bucket([]byte(collection))
		if bkt == nil {
			return nil
		}
		return bkt.Delete([]byte(key))
	})
}