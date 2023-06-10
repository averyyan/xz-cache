package xzcache

import (
	"fmt"
	"os"
	"testing"
	"time"

	xzmapitem "github.com/averyyan/xz-map/item"
)

var myCache Cache[string, int]

func TestMain(m *testing.M) {
	myCache = New[int](
		WithTickerInterval[string, int](1*time.Second),
		WithItemOpts[string, int](
			xzmapitem.WithDurationSeconds[int](2),
			xzmapitem.WithDeleteHandler[int](func(v int) {
				fmt.Println("this is deletehandler", v)
			}),
		),
	)
	for i := 0; i < 100; i++ {
		myCache.Set(fmt.Sprintf("%d", i), i)
	}
	os.Exit(m.Run())
}

func TestSet(t *testing.T) {
	myCache.Set(
		"aaaa",
		1000,
		xzmapitem.WithDurationSeconds[int](3),
		xzmapitem.WithDeleteHandler[int](func(v int) {
			fmt.Println("aaa is remove")
		}),
	)
	time.Sleep(5 * time.Second)
	v, ok := myCache.Get("1")
	t.Log(v, ok)

}
