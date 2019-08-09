// Volatile-timeout elimination, it is support for setting a cache item with expired seconds.
// When it reached capacity, some items will be deleted randomly in asynchronous queues.
package cache

import (
	"sync"
	"time"
)

const (
	CACHE_DEFAULT_TIMEOUT = 1800
	// randomly choose K caches to sample in asynchronous timer every time.
	RANDOM_CACHES_SELECT = 10
	// randomly choose C items to find minor expire item to delete.
	RANDOM_CANDIDATE_SIZE = 5
	// randomly choose K*M items to analyze per second.
	RANDOM_SAMPLES_SIZE = 100
	// if more than N items expired, then choose anthor M items to analyze.
	RANDOM_SAMPLES_TIMEOUT_SIZE = 25
)

type TCache struct {
	capacity int
	// the cache data.
	cacheMap map[string]*Item
	wlock    *sync.Mutex
	BasicCache
}

func NewTCache(capacity int) *TCache {
	return &TCache{
		capacity:   capacity,
		wlock:      new(sync.Mutex),
		cacheMap:   make(map[string]*Item, capacity),
		BasicCache: NewBasic(),
	}
}

func (this *TCache) GetMode() int {
	return MODE_TIMEOUT_VOLATILE
}

func (this *TCache) GetOrderSeq() uint64 {
	return 0
}

func (this *TCache) GetCapacity() int {
	return this.capacity
}

func (this *TCache) GetInfo() map[string]interface{} {
	res := make(map[string]interface{}, 4)
	res["length"] = len(this.cacheMap)
	res["capacity"] = this.GetCapacity()
	res["orderSeq"] = this.GetOrderSeq()
	res["mode"] = this.GetMode()
	return res
}

func (this *TCache) SetCapacity(capacity int) {
	this.capacity = capacity
}

func (this *TCache) Delete(key string) error {
	defer this.wlock.Unlock()
	this.wlock.Lock()
	if item, ok := this.cacheMap[key]; ok {
		delete(this.cacheMap, key)
		this.ReleaseItem(item)
	}
	return nil
}

func (this *TCache) Get(key string) interface{} {
	if item, ok := this.cacheMap[key]; ok {
		if item.seq < CurSecond {
			// this.Delete(key)
			return nil
		}
		return item.value
	}
	return nil
}

// Set a cache item with expired seconds.
func (this *TCache) SetEx(key string, value interface{}, timeout int) {
	duration := uint64(time.Now().Add(time.Duration(timeout) * time.Second).Unix())
	item := this.NewItem(duration, value)
	this.cacheMap[key] = item
}

func (this *TCache) Set(key string, value interface{}) {
	this.SetEx(key, value, CACHE_DEFAULT_TIMEOUT)
}

// Deleted items at random if reached capacity.
func (this *TCache) reduce() {
	// expire sooner (minor expire unix timestamp) is better candidate for deletion.
	i, j := 0, len(this.cacheMap)-this.capacity
	if j > 0 {
		bestKey, minSeq := "", uint64(0)
		for key, atom := range this.cacheMap {
			if i == 0 || atom.seq < minSeq {
				bestKey = key
				minSeq = atom.seq
			}
			i++

			if i >= RANDOM_CANDIDATE_SIZE {
				this.Delete(bestKey)
				j--
				i = 0
			}
		}
	}
}

// Remove expired items and do some cleaning up for capacity.
// On the limited sampling, we can decide whether or not to go further.
func (this *TCache) eliminate() {
	i, j, max := 0, 0, 2*RANDOM_SAMPLES_SIZE
	for key, item := range this.cacheMap {
		i++
		if item.seq < CurSecond {
			j++
			this.Delete(key)
		}
		if i > RANDOM_SAMPLES_SIZE {
			if j < RANDOM_SAMPLES_TIMEOUT_SIZE || i > max {
				break
			}
		}
	}

	this.reduce()
}

func (this *TCache) Check() int {
	if len(this.cacheMap) > RANDOM_SAMPLES_SIZE {
		this.eliminate()
		return 1
	}
	return 0
}
