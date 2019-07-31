package sample

import (
	"github.com/TianQinS/commhttp/cache"
	"github.com/kataras/iris"
	//"github.com/kataras/iris/context"
)

var (
	// specify alias can be used directly.
	Cache          *cache.Mgr
	NAME_NOT_EXIST map[string]interface{} = map[string]interface{}{
		"ok":   false,
		"data": "name not exist",
	}
)

// Client get cached data directly by set's name and data's key.
func GetItem(ctx iris.Context) {
	name := ctx.Params().Get("name")
	key := ctx.Params().Get("key")
	if r_cache, err := Cache.GetCache(name); err == nil {
		ctx.JSON(map[string]interface{}{
			"ok":   true,
			"data": r_cache.Get(key),
		})
		return
	}
	ctx.JSON(NAME_NOT_EXIST)
}

// Server set a data into a cache set,
// if the set is a volatile-timeout cache set, the timeout parameter will be provided.
func SetItem(ctx iris.Context) {
	name := ctx.Params().Get("name")
	key := ctx.Params().Get("key")
	if r_cache, err := Cache.GetCache(name); err == nil {
		value := ctx.PostValueDefault("value", "")
		timeout := ctx.PostValueIntDefault("timeout", 0)
		if value != "" && timeout >= 0 {
			r_cache.SetEx(key, value, timeout)
			ctx.JSON(map[string]interface{}{
				"ok":   true,
				"data": key,
			})
			return
		}
	}
	ctx.JSON(NAME_NOT_EXIST)
}

// Server delete a data from a cache set, no matter what type of the set.
func DeleteItem(ctx iris.Context) {
	name := ctx.Params().Get("name")
	key := ctx.Params().Get("key")
	if r_cache, err := Cache.GetCache(name); err == nil {
		r_cache.Delete(key)
		ctx.JSON(map[string]interface{}{
			"ok":   true,
			"data": key,
		})
		return
	}
	ctx.JSON(NAME_NOT_EXIST)
}

// Server init a cache set by name,
// the capacity and elimination algorithm type is required.
func InitCache(ctx iris.Context) {
	name := ctx.Params().Get("name")
	capacity, _ := ctx.Params().GetInt("capacity")
	mode, _ := ctx.Params().GetInt("mode")
	r_cache := Cache.InitCache(name, capacity, mode)
	ctx.JSON(map[string]interface{}{
		"ok":   true,
		"data": r_cache.GetInfo(),
	})
}

// Server delete a cache set by name.
func DeleteCache(ctx iris.Context) {
	name := ctx.Params().Get("name")
	ctx.JSON(map[string]interface{}{
		"ok":   true,
		"data": Cache.DeleteCache(name),
	})
}

func init() {
	Cache = cache.CacheMgr
}
