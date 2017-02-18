package server

import (
	"sync/atomic"
	"time"

	"gitlab.com/conspico/elasticshift/store"
)

// Cache ..
type Cache struct {
	ks  store.KeyStore
	key atomic.Value
}

func newCache(ks store.KeyStore) *Cache {
	return &Cache{ks: ks}
}

// GetKey ..
// Returns the cached key
func (c *Cache) GetKey() (store.Key, error) {

	key, ok := c.key.Load().(*store.Key)
	if ok && key != nil && time.Now().Before(key.NextSpin) {
		return *key, nil
	}

	nkey, err := c.ks.Get()
	if err != nil {
		return nkey, nil
	}

	if time.Now().Before(nkey.NextSpin) {
		c.key.Store(&nkey)
	}
	return nkey, nil
}
