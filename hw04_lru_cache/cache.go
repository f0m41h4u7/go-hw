package hw04_lru_cache //nolint:golint,stylecheck

import "sync"

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	capacity int
	queue    List
	items    map[Key]*listItem
	mutex    sync.Mutex
}

type cacheItem struct {
	key   Key
	value interface{}
}

func (cache *lruCache) Clear() {
	cache.mutex.Lock()
	defer cache.mutex.Unlock()

	// Clear map
	element := cache.queue.Back()
	for element != nil {
		delete(cache.items, element.Value.(cacheItem).key)
		element = element.Next
	}
	// Clear queue
	cache.queue = NewList()
}

func (cache *lruCache) Get(key Key) (interface{}, bool) {
	cache.mutex.Lock()
	defer cache.mutex.Unlock()

	// If key doesn't exist, return nil, false
	element, exists := cache.items[key]
	if !exists {
		return nil, false
	}

	// Move to front
	cache.queue.MoveToFront(element)
	return element.Value.(cacheItem).value, true
}

func (cache *lruCache) Set(key Key, value interface{}) bool {
	cache.mutex.Lock()
	defer cache.mutex.Unlock()

	// If key existed, return true
	element, exists := cache.items[key]
	newValue := cacheItem{
		key:   key,
		value: value,
	}

	if exists {
		element.Value = newValue
		cache.queue.MoveToFront(element)
		return true
	}

	// Delete last element if capacity is overflow
	if cache.queue.Len() == cache.capacity {
		lastElement := cache.queue.Back()
		cache.queue.Remove(lastElement)
		delete(cache.items, lastElement.Value.(cacheItem).key)
	}
	cache.queue.PushFront(newValue)
	cache.items[key] = cache.queue.Front()
	return false
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*listItem),
	}
}
