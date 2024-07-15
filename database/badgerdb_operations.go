package database

import (
	"blockchain-app/handlers"
	"errors"
	"fmt"
	"github.com/dgraph-io/badger/v4"
	"log"
	"os"
	"path/filepath"
	"strings"
)

const (
	dbPath = "./tmp/blocks_%s"
)

var dbs map[string]*badger.DB

func initDbs() {
	if dbs == nil {
		dbs = make(map[string]*badger.DB)
	}
}

func OpenDB(nodeID string) {
	initDbs()
	path := fmt.Sprintf(dbPath, nodeID)
	opts := badger.DefaultOptions(path)
	opts.Logger = nil
	db, err := openDatabase(path, opts)
	handlers.HandleErrors(err)
	dbs[nodeID] = db
}

func CloseDB(nodeID string) {
	if db, exists := dbs[nodeID]; exists {
		if err := db.Close(); err != nil {
			handlers.HandleErrors(err)
		}
		delete(dbs, nodeID)
	}
}

func Get(key []byte, nodeID string) []byte {
	db, exists := dbs[nodeID]
	if !exists {
		handlers.HandleErrors(errors.New("db closed"))
	}
	txn := db.NewTransaction(false)
	item, _ := txn.Get(key)
	if item == nil {
		return nil
	}

	byteItem, _ := item.ValueCopy(nil)
	return byteItem
}

func Set(key []byte, data []byte, nodeID string) {
	db, exists := dbs[nodeID]
	if !exists {
		handlers.HandleErrors(errors.New("db closed"))
	}
	txn := db.NewTransaction(true)
	setErr := txn.Set(key, data)
	handlers.HandleErrors(setErr)
	commitErr := txn.Commit()
	handlers.HandleErrors(commitErr)
}

func retry(dir string, originalOpts badger.Options) (*badger.DB, error) {
	lockPath := filepath.Join(dir, "LOCK")
	if err := os.Remove(lockPath); err != nil {
		return nil, fmt.Errorf(`removing "LOCK": %s`, err)
	}
	db, err := badger.Open(originalOpts)
	return db, err
}

func openDatabase(dir string, opts badger.Options) (*badger.DB, error) {
	if db, err := badger.Open(opts); err != nil {
		if strings.Contains(err.Error(), "lock") {
			if db, err := retry(dir, opts); err == nil {
				log.Println("database unlocked, value log truncated")
				return db, nil
			}
			log.Println("could not unlock database:", err)
		}
		return nil, err
	} else {
		return db, nil
	}
}
