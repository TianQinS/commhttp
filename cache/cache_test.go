package cache

import (
	//"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestVLru(t *testing.T) {
	r := CacheMgr.InitCache("vlru", 2, MODE_SLRU_VOLATILE)
	r.Set("1", "1")
	r.Set("2", "2")
	r.Get("2")
	r.Get("1")
	assert.Equal(t, "1", r.Get("1").(string))
	r.Set("3", "3")
	r.Get("3")
	assert.Equal(t, nil, r.Get("2"))
	assert.Equal(t, "1", r.Get("1").(string))
	r.Set("4", "4")
	r.Get("4")
	r.Get("1")
	assert.Equal(t, nil, r.Get("3"))
	assert.Equal(t, "1", r.Get("1").(string))
	r_cache, _ := CacheMgr.GetCache("vlru")
	assert.Equal(t, r_cache, r)
	CacheMgr.DeleteCache("vlru")
}

func TestTimeout(t *testing.T) {
	r := CacheMgr.InitCache("timeout", 2, MODE_TIMEOUT_VOLATILE)
	r.SetEx("1", "1", 2)
	time.Sleep(time.Second)
	assert.Equal(t, "1", r.Get("1").(string))
	time.Sleep(3 * time.Second)
	// for testing you should set up 'RANDOM_SAMPLES_SIZE' to 2.
	assert.Equal(t, nil, r.Get("1"))
	r_cache, _ := CacheMgr.GetCache("timeout")
	assert.Equal(t, r_cache, r)
	CacheMgr.DeleteCache("timeout")
}

func TestLru(t *testing.T) {
	r := CacheMgr.InitCache("lru", 2, MODE_LRU_ORIGINAL)
	r.Set("1", "1")
	r.Set("2", "2")
	r.Set("3", "3")
	assert.Equal(t, "2", r.Get("2").(string))
	assert.Equal(t, nil, r.Get("1"))
	r.Set("4", "4")
	assert.Equal(t, nil, r.Get("3"))
	r_cache, _ := CacheMgr.GetCache("lru")
	assert.Equal(t, r_cache, r)
	CacheMgr.DeleteCache("lru")
}
