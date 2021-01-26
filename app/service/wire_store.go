package service

import (
	"github.com/suisrc/zgo/modules/store"
	"github.com/suisrc/zgo/modules/store/buntdb"
)

// NewStorer 全局缓存
func NewStorer() (store.Storer, func(), error) {
	store, err := buntdb.NewStore(":memory:") // 使用内存缓存
	if err != nil {
		return nil, nil, err
	}
	return store, func() { store.Close() }, nil
}
