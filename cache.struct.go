package xzcache

import (
	"sync"

	xzmapitem "github.com/averyyan/xz-map/item"
	xzmap "github.com/averyyan/xz-map/map"
	xzticker "github.com/averyyan/xz-ticker"
)

type Cache[K comparable, T any] interface {
	Has(key K) bool
	Get(key K) (T, bool)
	Set(key K, value T, opts ...xzmapitem.Option[T])
	Remove(key K)
	Values() map[K]T
	IterBuffered() <-chan cacheTuple[K, T]
}

type cache[K comparable, T any] struct {
	data       *xzmap.Map[K, T]     // 数据储存
	ticker     xzticker.Ticker      // 定时器
	mapOps     []xzmap.Option[K, T] // map options
	tickerOpts []xzticker.Option    // 定时器 options
	itemOps    []xzmapitem.Option[T]
}

// 是否存在缓存
func (c *cache[K, T]) Has(key K) bool {
	return c.data.Has(key)
}

// 获取缓存值
func (c *cache[K, T]) Get(key K) (T, bool) {
	item, ok := c.data.Get(key)
	var value T
	if item == nil {
		return value, ok
	}
	return item.GetValue(), ok
}

// 设置缓存
func (c *cache[K, T]) Set(key K, value T, opts ...xzmapitem.Option[T]) {
	dOpts := c.itemOps
	dOpts = append(dOpts, opts...)
	c.data.Set(key, value, dOpts...)
}

// 删除缓存
func (c *cache[K, T]) Remove(key K) {
	c.data.Remove(key)
}

// 获取所有的缓存值
func (c *cache[K, T]) Values() map[K]T {
	items := make(map[K]T)
	for key, item := range c.data.Items() {
		items[key] = item.GetValue()
	}
	return items
}

// 同步 缓存清理
func (c *cache[K, T]) Clean() {
	for key := range c.Values() {
		c.data.Remove(key)
	}
}

type cacheTuple[K comparable, T any] struct {
	Key   K
	Value T
}

// 管道迭代
func (c *cache[K, T]) IterBuffered() <-chan cacheTuple[K, T] {
	chans := make(chan cacheTuple[K, T])
	wg := sync.WaitGroup{}
	wg.Add(c.data.Size())
	go func() {
		for tuple := range c.data.IterBuffered() {
			chans <- cacheTuple[K, T]{Key: tuple.Key, Value: tuple.Val.GetValue()}
			wg.Done()
		}
		wg.Wait()
		close(chans)
	}()
	return chans
}

// 缓存停止删除
func (c *cache[K, T]) Stop() {
	c.ticker.Stop()
}
