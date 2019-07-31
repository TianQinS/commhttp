// volatile-lru
package cache

import (
	"sync"
)

const (
	// randomly choose some items to analyze.
	SAMPLES_VOLATILE_SIZE = 5
	// filter whether we can check the item.
	SAMPLES_FILTER_RATIO = 0.6
	SEQ_LARGEST          = uint64(2 ^ 64 - 1)
	// dealing with the out-of-bounds, even though this rarely happens.
	SEQ_THRESHOLD = uint64(2 ^ 64 - 2 ^ 32)
)

type VLRUCache struct {
	capacity   int
	checkPoint uint64
	// the recent visit seq.
	orderSeq uint64
	// the cache data.
	cacheMap map[string]*Item
	wlock    *sync.Mutex
	BasicCache
}

func NewVLRUCache(capacity int) *VLRUCache {
	return &VLRUCache{
		capacity:   capacity,
		checkPoint: uint64(float32(capacity) * SAMPLES_FILTER_RATIO),
		orderSeq:   0,
		wlock:      new(sync.Mutex),
		cacheMap:   make(map[string]*Item, capacity),
		BasicCache: NewBasic(),
	}
}

func (this *VLRUCache) GetMode() int {
	return MODE_SLRU_VOLATILE
}

func (this *VLRUCache) GetOrderSeq() uint64 {
	return this.orderSeq
}

func (this *VLRUCache) GetCapacity() int {
	return this.capacity
}

func (this *VLRUCache) GetInfo() map[string]interface{} {
	res := make(map[string]interface{}, 4)
	res["length"] = len(this.cacheMap)
	res["capacity"] = this.GetCapacity()
	res["orderSeq"] = this.GetOrderSeq()
	res["mode"] = this.GetMode()
	return res
}

func (this *VLRUCache) SetCapacity(capacity int) {
	this.capacity = capacity
	this.checkPoint = uint64(float32(capacity) * SAMPLES_FILTER_RATIO)
}

func (this *VLRUCache) Delete(key string) error {
	defer this.wlock.Unlock()
	this.wlock.Lock()
	if item, ok := this.cacheMap[key]; ok {
		delete(this.cacheMap, key)
		this.ReleaseItem(item)
	}
	return nil
}

// update most recently used seq,
func (this *VLRUCache) Get(key string) interface{} {
	if item, ok := this.cacheMap[key]; ok {
		// atomic.AddUint64(&this.orderSeq, 1)
		this.orderSeq++
		item.seq = this.orderSeq
		return item.value
	}
	return nil
}

func (this *VLRUCache) Set(key string, value interface{}) {
	item := this.NewItem(this.orderSeq, value)
	this.cacheMap[key] = item
	if len(this.cacheMap) > this.capacity {
		this.eliminate()
	}
}

func (this *VLRUCache) SetEx(key string, value interface{}, timeout int) {
	this.Set(key, value)
}

// volatile-lru elimination without eviction pool.
func (this *VLRUCache) eliminate() {
	oldest_key := ""
	cnt, cur, oldest := 0, uint64(0), SEQ_LARGEST
	point := this.orderSeq - this.checkPoint

	// the sample is chosed by the range random.
	for key, atom := range this.cacheMap {
		cur = atom.seq
		if cur <= point {
			cnt++
			if oldest > cur {
				oldest = cur
				oldest_key = key
			}
			if cnt >= SAMPLES_VOLATILE_SIZE {
				break
			}
		}
	}
	if oldest_key != "" {
		this.Delete(oldest_key)
	}
}

func (this *VLRUCache) Check() int {
	if this.orderSeq > SEQ_THRESHOLD {
		return 1
	}
	return 0
}
