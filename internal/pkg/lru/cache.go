package lru

import (
	"sync"
)

// Key is the type of keys used in the cache.
type Key = string

// ICache is an interface that represents a simple key-value store with LRU eviction policy.
type ICache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	mutex    *sync.Mutex
	capacity int
	queue    List
	items    map[Key]*ListItem
}

// NewCache creates a new instance of the LRU cache with the specified capacity.
func NewCache(capacity int) ICache {
	return &lruCache{
		mutex:    new(sync.Mutex),
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}

type cachedValue struct {
	key   Key
	value any
}

// Set adds or updates a key-value pair in the cache. If the key already exists,
// it updates the value and moves it to the front of the queue.
// If the key does not exist, it adds a new item to the front of the queue.
func (lc *lruCache) Set(key Key, value any) bool {
	lc.mutex.Lock()
	defer lc.mutex.Unlock()

	cv := cachedValue{key: key, value: value}

	item, ok := lc.items[key]
	if ok {
		item.Value = cv
		lc.queue.MoveToFront(item)
		return true
	}

	lc.items[key] = lc.queue.PushFront(cv)

	if lc.queue.Len() > lc.capacity {
		last := lc.queue.Back()
		lc.queue.Remove(last)
		delete(lc.items, last.Value.(cachedValue).key)
	}

	return false
}

// Get retrieves the value associated with a given key from the cache.
func (lc *lruCache) Get(key Key) (any, bool) {
	lc.mutex.Lock()
	defer lc.mutex.Unlock()

	if item, ok := lc.items[key]; ok {
		lc.queue.MoveToFront(item)
		return item.Value.(cachedValue).value, true
	}

	return nil, false
}

// Clear removes all items from the cache.
func (lc *lruCache) Clear() {
	lc.mutex.Lock()
	defer lc.mutex.Unlock()

	lc.items = make(map[Key]*ListItem, lc.capacity)
	lc.queue = NewList()
}
