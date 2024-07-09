package main

import (
	"fmt"
	"memCache/cache-server"
	"time"
)

func main(){
	cache := cache_server.NewMemCache()
	cache.SetMaxMemory("100MB")
	cache.Set("int", 1, time.Second)
	cache.Set("bool", false, time.Second)
	cache.Set("data", map[string]interface{}{"a": 1}, time.Second)
	fmt.Println(cache.Get("int"))
	fmt.Println(cache.Get("data"))
	// fmt.Println(cache.Keys())
	// time.Sleep(time.Second*3)
	fmt.Println(cache.Keys())
	cache.Del("int")
	fmt.Println(cache.Keys())
	cache.Flush()
	fmt.Println(cache.Keys())
}