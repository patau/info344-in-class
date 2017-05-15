package main

import (
	"fmt"
	"sync"
	"time"
)

//CacheEntry represents an entry in the cache
type CacheEntry struct {
	value   string
	expires time.Time
}

//Cache represents a map[string]string that is safe
//for concurrent access
type Cache struct {
	//TODO: protect this map with a RWMutex
	mu      sync.RWMutex
	entries map[string]*CacheEntry
	quit    chan bool
}

//NewCache creates and returns a new Cache
func NewCache() *Cache {
	c := &Cache{
		entries: make(map[string]*CacheEntry),
		mu:      sync.RWMutex{},
		quit:    make(chan bool),
	}
	go c.startJanitor() //go ensures that startJ will run on its own?
	return c
}

//go routine will end and allow for cleaning up correctly
func (c *Cache) Close() {
	c.quit <- true //writes true value to c.quit
}

func (c *Cache) startJanitor() {
	ticker := time.NewTicker(time.Second) //Channel; every second, something written to chan
	for {
		//tries to do these channel processes; does whatever is ready
		select {
		case <-ticker.C: //read from ticker.Cache
			//if ready to be read from
			c.purgeExpired() //called every time ticker
		case <-c.quit: //release the cache - allow it to be garbage collected
			return
		}
	}
}

//get rid of anything expired in the cache
func (c *Cache) purgeExpired() {
	c.mu.Lock()
	defer c.mu.Unlock()

	now := time.Now()
	nPurged := 0 //increment everytime we delete something
	for key, entry := range c.entries {
		if now.After(entry.expires) {
			//time to remove from the map b/c expired
			delete(c.entries, key)
			nPurged++
		}
	}
	fmt.Printf("purged %d entries\n", nPurged) //just for dev purposes
}

//Get returns the value associated with the requested key.
//The returned boolean will be false if the key was not
//in the cache.
func (c *Cache) Get(key string) (string, bool) {
	c.mu.RLock() //get a read lock - need to make sure it gets released no matter what!
	defer c.mu.RUnlock()
	entry := c.entries[key]
	if entry == nil {
		return "", false
	}
	return entry.value, true
}

//Set sets the value associated with the given key.
//If the key is not yet in the cache, it will be added.
func (c *Cache) Set(key string, value string, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	entry := c.entries[key]
	if entry == nil {
		entry = &CacheEntry{} //if the item is not already in the map, add new for that key
		c.entries[key] = entry
	}
	entry.value = value
	entry.expires = time.Now().Add(ttl)
}
