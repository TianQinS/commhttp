// The volatile elimination algorithm for cache.
package cache

import (
	"fmt"
	"sync"
	"time"

	"github.com/TianQinS/fastapi/timer"
)

const (
	// volatile-lru policy
	MODE_SLRU_VOLATILE = iota
	// volatile-timeout policy
	MODE_TIMEOUT_VOLATILE
	// original lru
	MODE_LRU_ORIGINAL

	// default check interval for timeout policy.
	CACHE_CHECK_INTERVAL = time.Second
	// shared for concurrency conflicts
	CACHE_SHARED_LENGTH = 10
)

var (
	CurSecond uint64
	// keeps all caches.
	CacheMgr *Mgr
)

type (
	Item struct {
		seq   uint64
		value interface{}
	}

	BasicCache struct {
		lock *sync.RWMutex
		// for read and delete are in concurrency conflict.
		itemShare []*Item
		itemPool  sync.Pool
	}

	Mgr struct {
		cacheMgr map[string]Cache
		lock     *sync.Mutex
	}
)

// Cache defines the interface of elimination algorithm object.
type Cache interface {
	GetMode() int
	GetOrderSeq() uint64
	GetCapacity() int
	SetCapacity(capacity int)
	Check() int

	GetInfo() map[string]interface{}
	Set(key string, value interface{})
	SetEx(key string, value interface{}, timeout int)
	// return nil when missed.
	Get(key string) interface{}
	Delete(key string) error
}

// BasicCache's default constructor.
func NewBasic() BasicCache {
	return BasicCache{
		lock:      new(sync.RWMutex),
		itemShare: make([]*Item, 0, CACHE_SHARED_LENGTH),
		itemPool: sync.Pool{
			New: func() interface{} {
				return &Item{
					seq:   0,
					value: nil,
				}
			},
		},
	}
}

// Get a pooled item.
func (this *BasicCache) NewItem(seq uint64, value interface{}) *Item {
	item := this.itemPool.Get().(*Item)
	item.seq = seq
	item.value = value
	return item
}

// Put an item into the pool.
func (this *BasicCache) PutItem() {
	this.lock.Lock()
	oldest := this.itemShare[0]
	this.itemShare = this.itemShare[1:]
	this.lock.Unlock()
	this.itemPool.Put(oldest)
}

// Append an item to the shared queue.
func (this *BasicCache) ReleaseItem(item *Item) {
	if item != nil {
		this.lock.Lock()
		this.itemShare = append(this.itemShare, item)
		this.lock.Unlock()
		length := len(this.itemShare)
		if length >= CACHE_SHARED_LENGTH {
			this.PutItem()
		}
	}
}

func (this *Mgr) GetCache(key string) (Cache, error) {
	if cache, ok := this.cacheMgr[key]; ok {
		return cache, nil
	}
	return nil, fmt.Errorf("cache=%s not exist", key)
}

func (this *Mgr) DeleteCache(key string) error {
	defer this.lock.Unlock()
	this.lock.Lock()
	if _, ok := this.cacheMgr[key]; ok {
		delete(this.cacheMgr, key)
		return nil
	} else {
		return fmt.Errorf("cache=%s not exist", key)
	}
}

func (this *Mgr) InitCache(key string, capacity, mode int) Cache {
	defer this.lock.Unlock()
	this.lock.Lock()
	var cache Cache
	if cache, ok := this.cacheMgr[key]; ok {
		cache.SetCapacity(capacity)
		return cache
	}

	switch mode {
	case MODE_SLRU_VOLATILE:
		cache = NewVLRUCache(capacity)
	case MODE_TIMEOUT_VOLATILE:
		cache = NewTCache(capacity)
	case MODE_LRU_ORIGINAL:
		cache = NewLRUCache(capacity)
	default:
		cache = NewTCache(capacity)
	}
	this.cacheMgr[key] = cache
	return cache
}

func (this *Mgr) Check() {
	cnt := 0
	// update current time.
	CurSecond = uint64(time.Now().Unix())
	for key, cache := range this.cacheMgr {
		mode := cache.GetMode()
		// for timeout cache.
		if mode == MODE_TIMEOUT_VOLATILE {
			if cnt >= RANDOM_CACHES_SELECT {
				continue
			}
			cnt += cache.Check()
		} else if mode == MODE_SLRU_VOLATILE {
			if cache.Check() == 1 {
				this.DeleteCache(key)
				this.InitCache(key, cache.GetCapacity(), MODE_SLRU_VOLATILE)
			}
		}
	}
}

func init() {
	CacheMgr = &Mgr{
		lock:     new(sync.Mutex),
		cacheMgr: make(map[string]Cache, 0),
	}

	timer.AddTimer(CACHE_CHECK_INTERVAL, CacheMgr.Check)
}
