package cache

import (
	"log"
	"sync"
	"time"
)

type memCache struct {
	// 最大内存
	maxMemorySize int64
	// 最大内存字符串表示
	maxMemorySizeStr string
	// 当前已使用内存
	currentMemorySize int64
	// 存储键值对
	values map[string]*memCacheValue
	// 读写锁
	locker sync.RWMutex
	// 清楚过期缓存时间间隔
	clearExpiredItemTimeInterval time.Duration
}

type memCacheValue struct{
	// value值
	val interface{}
	// 过期时间
	expireTime time.Time
	// 有效时长
	expire time.Duration
	// value大小
	size int64
}

func NewMemCache() Cache{
	mc := &memCache{
		currentMemorySize: 0,
		values: make(map[string]*memCacheValue),
		clearExpiredItemTimeInterval: time.Second,
	}
	go mc.clearExpiredItem()
	return mc
}

// size: 1KB 100KB 1MB 2MB 1GB
func (mc *memCache) SetMaxMemory(size string) bool {
	// fmt.Println("called set max memory")
	mc.maxMemorySize, mc.maxMemorySizeStr = ParseSize(size)
	return true
}

// 将value写入缓存
func (mc *memCache) Set(key string, val interface{}, expire time.Duration) bool {
	// fmt.Println("called set")
	mc.locker.Lock()
	defer mc.locker.Unlock()
	v := &memCacheValue{
		val: val,
		expireTime: time.Now().Add(expire),
		expire: expire,
		size: GetValSize(val),
	}
	mc.del(key)
	mc.add(key, v)
	if(mc.currentMemorySize>mc.maxMemorySize){
		mc.del(key)
		log.Printf("超出最大内存限制%s",mc.maxMemorySizeStr)
		//panic(fmt.Sprintf("超出最大内存限制%s",mc.maxMemorySizeStr))
	}
	return true
}

func (mc *memCache) get(key string) (*memCacheValue, bool){
	val, ok := mc.values[key]
	return val, ok
}

func (mc *memCache) del(key string){
	tmp, ok := mc.get(key)
	if ok && tmp!=nil {
		mc.currentMemorySize -= tmp.size
		delete(mc.values, key)
	}
}

func (mc *memCache) add(key string, val *memCacheValue){
	mc.values[key] = val 
	mc.currentMemorySize += val.size
}

// 根据key值获取value
func (mc *memCache) Get(key string) (interface{}, bool) {
	// fmt.Println("called get")
	mc.locker.RLock()
	defer mc.locker.RUnlock()
	mcv, ok := mc.get(key)
	if ok {
		// 判断缓存是否过期
		if mcv.expire!=0&&mcv.expireTime.Before(time.Now()){
			mc.del(key)
			return nil, false
		}
		return mcv.val, ok
	}
	return nil, false
}

// 删除key值
func (mc *memCache) Del(key string) bool {
	// fmt.Println("called del")
	mc.locker.Lock()
	defer mc.locker.Unlock()
	mc.del(key)
	return true
}

// 判断key值是否存在
func (mc *memCache) Exists(key string) bool {
	// fmt.Println("called exists")
	mc.locker.RLock()
	defer mc.locker.RUnlock()
	_, ok := mc.values[key]
	return ok
}

// 清空所有的key
func (mc *memCache) Flush() bool {
	// fmt.Println("called flush")
	mc.locker.Lock()
	defer mc.locker.Unlock()
	mc.values = make(map[string]*memCacheValue, 0)
	mc.currentMemorySize = 0
	return true
}

// 获取缓存中所有key的数量
func (mc *memCache) Keys() int64 {
	// fmt.Println("called keys")
	mc.locker.RLock()
	defer mc.locker.RUnlock()
	return int64(len(mc.values))
}

func (mc *memCache) clearExpiredItem(){
	timeTicker := time.NewTicker(mc.clearExpiredItemTimeInterval)
	defer timeTicker.Stop()
	for range timeTicker.C{
		for key, item := range mc.values{
			if item.expire!=0&&item.expireTime.Before(time.Now()){
				mc.locker.Lock()
				mc.del(key)
				mc.locker.Unlock()
			}			
		}
	}
}