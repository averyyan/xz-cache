package xzcache

import (
	"time"

	xzmap "github.com/averyyan/xz-map/map"
	xzticker "github.com/averyyan/xz-ticker"
)

func New[T any](opts ...Option[string, T]) Cache[string, T] {
	c := &cache[string, T]{
		mapOps:     make([]xzmap.Option[string, T], 0),
		tickerOpts: make([]xzticker.Option, 0),
	}
	// 注入定时数据清理函数
	opts = append(opts, WithTickerHandler[string, T](func() {
		deleteExpired(c)
	}))
	for _, opt := range opts {
		opt(c)
	}
	// 初始化数据
	c.data = xzmap.New[T](c.mapOps...)
	// 初始化定时器
	c.ticker = xzticker.New(c.tickerOpts...)
	go c.ticker.Run()
	return c
}

func deleteExpired[T any](c *cache[string, T]) {
	now := time.Now().UnixNano()
	for tuple := range c.data.IterBuffered() {
		if tuple.Val.VerifyExpiration(now) {
			c.data.Remove(tuple.Key)
		}
	}
}
