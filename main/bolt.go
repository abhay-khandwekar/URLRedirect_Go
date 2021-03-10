package main

import (
	"fmt"
	"log"

	"github.com/boltdb/bolt"
)

var dBase *bolt.DB

func init() {
	db, err := bolt.Open("pathURL.db", 0600, nil)
	if err != nil {
		log.Fatal("Could not initialize Bolt-db: pathURL.db.")
	}

	err = db.Update(func(tx *bolt.Tx) error {
		pathurlBkt, err := tx.CreateBucketIfNotExists([]byte("PATH_URL"))
		if err != nil {
			return fmt.Errorf("Could not create bucket PATH_URL")
		}

		err = pathurlBkt.Put([]byte("/gojson"), []byte("https://golang.org/pkg/encoding/json/"))
		if err != nil {
			return fmt.Errorf("Could not create path: gojson")
		}

		err = pathurlBkt.Put([]byte("/bolt"), []byte("https://github.com/boltdb/bolt"))
		if err != nil {
			return fmt.Errorf("Could not create path: gojson")
		}

		dBase = db
		fmt.Println("Db initialized.")

		return nil
	})
}
