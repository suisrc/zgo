package buntdb

import (
	"context"
	"os"
	"path/filepath"
	"time"

	"github.com/suisrc/zgo/store"
	"github.com/tidwall/buntdb"
)

// NewStore 创建基于buntdb的文件存储
func NewStore(path string) (*Store, error) {
	if path != ":memory:" {
		os.MkdirAll(filepath.Dir(path), 0777)
	}

	db, err := buntdb.Open(path)
	if err != nil {
		return nil, err
	}

	return &Store{
		db: db,
	}, nil
}

var _ store.Storer = new(Store)

// Store buntdb存储
type Store struct {
	db *buntdb.DB
}

// TTL ...
func (a *Store) TTL(ctx context.Context, key string) (time.Duration, bool, error) {
	var exists bool
	var value time.Duration
	err := a.db.View(func(tx *buntdb.Tx) error {
		val, err := tx.TTL(key)
		if err != nil {
			if err == buntdb.ErrNotFound {
				return nil
			}
			return err
		}
		value = val
		exists = true
		return nil
	})
	return value, exists, err
}

// EXP ...
func (a *Store) EXP(ctx context.Context, key string, expiration time.Duration) (bool, error) {
	var exists bool
	err := a.db.Update(func(tx *buntdb.Tx) error {
		val, err := tx.Get(key)
		if err != nil {
			if err == buntdb.ErrNotFound {
				return nil // 没有数据
			}
			return err
		}

		var opts *buntdb.SetOptions
		if expiration > 0 {
			opts = &buntdb.SetOptions{Expires: true, TTL: expiration}
		}
		_, _, err = tx.Set(key, val, opts)
		exists = true
		return err
	})
	return exists, err
}

// Set ...
func (a *Store) Set(ctx context.Context, key, value string, expiration time.Duration) error {
	return a.db.Update(func(tx *buntdb.Tx) error {
		var opts *buntdb.SetOptions
		if expiration > 0 {
			opts = &buntdb.SetOptions{Expires: true, TTL: expiration}
		}
		_, _, err := tx.Set(key, value, opts)
		return err
	})
}

// Get ...
func (a *Store) Get(ctx context.Context, key string) (string, bool, error) {
	var exists bool
	var value string
	err := a.db.View(func(tx *buntdb.Tx) error {
		val, err := tx.Get(key)
		if err != nil {
			if err == buntdb.ErrNotFound {
				return nil
			}
			return err
		}
		value = val
		exists = true
		return nil
	})
	return value, exists, err
}

// Set1 ...
func (a *Store) Set1(ctx context.Context, key string, expiration time.Duration) error {
	return a.db.Update(func(tx *buntdb.Tx) error {
		var opts *buntdb.SetOptions
		if expiration > 0 {
			opts = &buntdb.SetOptions{Expires: true, TTL: expiration}
		}
		_, _, err := tx.Set(key, "1", opts)
		return err
	})
}

// Delete ...
func (a *Store) Delete(ctx context.Context, key string) error {
	return a.db.Update(func(tx *buntdb.Tx) error {
		_, err := tx.Delete(key)
		if err != nil && err != buntdb.ErrNotFound {
			return err
		}
		return nil
	})
}

// Check ...
func (a *Store) Check(ctx context.Context, key string) (bool, error) {
	var exists bool
	err := a.db.View(func(tx *buntdb.Tx) error {
		val, err := tx.Get(key)
		if err != nil && err != buntdb.ErrNotFound {
			return err
		}
		exists = val != ""
		return nil
	})
	return exists, err
}

// Close ...
func (a *Store) Close() error {
	return a.db.Close()
}
