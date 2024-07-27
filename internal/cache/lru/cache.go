package lru

import (
	"container/list"
	"sync"
)

var _ Cache = (*CacheLRU)(nil)

type Cache interface {
	Set(key string, item *Item) bool
	Get(key string) (*Item, bool)
}

type Item struct {
	key      string
	FileName string
	Size     uint64
}

type RemoveItemCallback func(item *Item)

type CacheLRU struct {
	mu           sync.RWMutex
	list         *list.List
	items        map[string]*list.Element
	limit        uint64
	size         uint64
	onRemoveFunc RemoveItemCallback
}

func CreateCacheLRU(limit uint64, onRemoveFunc RemoveItemCallback) *CacheLRU {
	return &CacheLRU{
		list:         list.New(),
		items:        make(map[string]*list.Element),
		limit:        limit,
		onRemoveFunc: onRemoveFunc,
	}
}

func (c *CacheLRU) Set(key string, item *Item) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	item.key = key
	listItem, isExists := c.items[key]
	if isExists {
		cacheItem := listItem.Value.(*list.Element)
		cacheItem.Value = item
		c.list.MoveToFront(listItem)

		return isExists
	}

	if c.isCapacityExceeded() {
		c.dropLastItem()
	}

	newListItem := c.list.PushFront(item)
	c.items[key] = newListItem
	c.size += item.Size

	return isExists
}

func (c *CacheLRU) Get(key string) (*Item, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	listItem, isExists := c.items[key]
	if !isExists {
		return nil, isExists
	}

	c.list.MoveToFront(listItem)

	return listItem.Value.(*Item), isExists
}

func (c *CacheLRU) dropLastItem() {
	lastListItem := c.list.Back()
	item := lastListItem.Value.(*Item)
	element, isExists := c.items[item.key]
	if !isExists {
		return
	}

	delete(c.items, item.key)
	c.list.Remove(element)
	c.size -= item.Size

	c.onRemoveFunc(item)
}

func (c *CacheLRU) isCapacityExceeded() bool {
	return c.size >= c.limit
}
