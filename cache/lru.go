// original lru
package cache

import (
	"container/list"
	"sync"
)

type ItemLRU struct {
	key   string
	value interface{}
}

type LRUCache struct {
	capacity int
	dList    *list.List
	// the cache data.
	cacheMap map[string]*list.Element
	lock     *sync.Mutex
}

func NewLRUCache(capacity int) *LRUCache {
	return &LRUCache{
		capacity: capacity,
		dList:    list.New(),
		lock:     new(sync.Mutex),
		cacheMap: make(map[string]*list.Element, capacity),
	}
}

func (this *LRUCache) GetMode() int {
	return MODE_LRU_ORIGINAL
}

func (this *LRUCache) GetOrderSeq() uint64 {
	return uint64(0)
}

func (this *LRUCache) GetCapacity() int {
	return this.capacity
}

func (this *LRUCache) GetInfo() map[string]interface{} {
	res := make(map[string]interface{}, 4)
	res["length"] = len(this.cacheMap)
	res["capacity"] = this.GetCapacity()
	res["orderSeq"] = this.GetOrderSeq()
	res["mode"] = this.GetMode()
	return res
}

func (this *LRUCache) SetCapacity(capacity int) {
	this.capacity = capacity
}

func (this *LRUCache) Delete(key string) error {
	defer this.lock.Unlock()
	this.lock.Lock()
	if node, ok := this.cacheMap[key]; ok {
		item := node.Value.(*ItemLRU)
		delete(this.cacheMap, item.key)
		this.dList.Remove(node)
	}
	return nil
}

func (this *LRUCache) Get(key string) interface{} {
	defer this.lock.Unlock()
	this.lock.Lock()
	if node, ok := this.cacheMap[key]; ok {
		item := node.Value.(*ItemLRU)
		this.dList.MoveToFront(node)
		return item.value
	}
	return nil
}

func (this *LRUCache) Set(key string, value interface{}) {
	defer this.lock.Unlock()
	this.lock.Lock()
	if node, ok := this.cacheMap[key]; ok {
		this.dList.MoveToFront(node)
		node.Value.(*ItemLRU).value = value
		return
	}

	newNode := this.dList.PushFront(&ItemLRU{
		key:   key,
		value: value,
	})
	this.cacheMap[key] = newNode
	if this.dList.Len() > this.capacity {
		this.eliminate()
	}
}

func (this *LRUCache) SetEx(key string, value interface{}, timeout int) {
	this.Set(key, value)
}

// Remove the last.
func (this *LRUCache) eliminate() {
	node := this.dList.Back()
	if node == nil {
		return
	}
	item := node.Value.(*ItemLRU)
	delete(this.cacheMap, item.key)
	this.dList.Remove(node)
}

func (this *LRUCache) Check() int {
	return 0
}
