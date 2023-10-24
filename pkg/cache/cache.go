package cache

import (
	"container/list"
	"errors"
	"fmt"
	"sync"
)

// Cache implements a Least Frequently Used cache where
// if storage reaches max capacity, the least used items are evicted.
type Cache struct {
	sync.RWMutex
	storage     map[any]*list.Element
	elemByUsage map[int]*list.List
	size        int
	minFreq     int
}

type item struct {
	key       any
	val       any
	usageFreq int
}

// NewCache creates and returns a new Cache object
// with storage size for up to maxItems elements.
func NewCache(maxItems int) *Cache {
	return &Cache{
		storage:     make(map[any]*list.Element),
		elemByUsage: make(map[int]*list.List),
		size:        maxItems,
		minFreq:     1,
	}
}

var ErrKeyNotFound = errors.New("key not found in cache")

// Get retrieves the value of the provided key if it exists.
// If the key does not exist, it returns an error.
func (c *Cache) Get(key any) (any, error) {
	c.RLock()
	defer c.RUnlock()

	found, ok := c.storage[key]
	if !ok {
		return nil, fmt.Errorf("%s: %w", key, ErrKeyNotFound)
	}

	currItem := found.Value.(*item)
	c.updateFreq(key, currItem, found)

	return currItem.val, nil
}

// updateFreq updates the frequency count of a cache item.
func (c *Cache) updateFreq(key any, currItem *item, currElem *list.Element) {
	c.elemByUsage[currItem.usageFreq].Remove(currElem)
	currItem.usageFreq++

	itemsToUpdate, ok := c.elemByUsage[currItem.usageFreq]
	if !ok {
		itemsToUpdate = list.New()
	}
	updatedItem := itemsToUpdate.PushFront(currItem)

	c.elemByUsage[currItem.usageFreq] = itemsToUpdate
	c.storage[key] = updatedItem

	// Update minFreq if no other items are left.
	if currItem.usageFreq-1 == c.minFreq && c.elemByUsage[c.minFreq].Len() == 0 {
		c.minFreq++
	}
}

// Set adds a new key-value pair to the cache or updates the existing one.
func (c *Cache) Set(key, val any) {
	if c.size == 0 {
		return
	}

	c.Lock()
	defer c.Unlock()

	// if exist — update (increase) frequency for key
	if elem, exists := c.storage[key]; exists {
		currItem := elem.Value.(*item)
		currItem.val = val
		c.updateFreq(key, currItem, elem)
		return
	}

	// if storage reached limit — remove one item
	if len(c.storage) == c.size {
		c.evict()
	}

	c.minFreq = 1
	currItem := &item{
		key:       key,
		val:       val,
		usageFreq: 1,
	}
	if _, exists := c.elemByUsage[1]; !exists {
		c.elemByUsage[1] = list.New()
	}
	c.storage[key] = c.elemByUsage[1].PushFront(currItem)
}

// evict removes the least frequently used item from the cache.
func (c *Cache) evict() {
	minList := c.elemByUsage[c.minFreq]
	backNode := minList.Back()
	if backNode == nil {
		return
	}
	delete(c.storage, backNode.Value.(*item).key)
	minList.Remove(backNode)
	// clean up empty minList
	if minList.Len() == 0 {
		delete(c.elemByUsage, c.minFreq)
		c.minFreq++
	}
}
