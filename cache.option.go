package xzcache

import (
	"time"

	mapitem "github.com/averyyan/xz-map/item"
	xzmap "github.com/averyyan/xz-map/map"
	xzticker "github.com/averyyan/xz-ticker"
)

type Option[K comparable, T any] func(c *cache[K, T])

// 配置定时间隔
func WithTickerInterval[K comparable, T any](interval time.Duration) Option[K, T] {
	return func(c *cache[K, T]) {
		c.tickerOpts = append(c.tickerOpts, xzticker.WithInterval(interval))
	}
}

// 配置缓存处理
func WithTickerHandler[K comparable, T any](handler func()) Option[K, T] {
	return func(c *cache[K, T]) {
		c.tickerOpts = append(c.tickerOpts, xzticker.WithHandler(handler))
	}
}

// 配置map分片数量
func WithMapSharedSize[K comparable, T any](size int) Option[K, T] {
	return func(c *cache[K, T]) {
		c.mapOps = append(c.mapOps, xzmap.WithSharedSize[K, T](size))
	}
}

// 设置默认 map item 配置
func WithItemOpts[K comparable, T any](opts ...mapitem.Option[T]) Option[K, T] {
	return func(c *cache[K, T]) {
		c.itemOps = append(c.itemOps, opts...)
	}
}
