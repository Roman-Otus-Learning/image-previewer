package lru

import (
	"math/rand"
	"strconv"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCache(t *testing.T) {
	t.Run("empty cache", func(t *testing.T) {
		c := CreateCacheLRU(100, func(item *Item) {})

		_, ok := c.Get("aaa")
		require.False(t, ok)

		_, ok = c.Get("bbb")
		require.False(t, ok)
	})

	t.Run("simple", func(t *testing.T) {
		c := CreateCacheLRU(100, func(item *Item) {})

		wasInCache := c.Set("aaa", createItemStub("aaa", "filename1", 10))
		require.False(t, wasInCache)

		wasInCache = c.Set("bbb", createItemStub("bbb", "filename2", 10))
		require.False(t, wasInCache)

		val, ok := c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, "filename1", val.FileName)

		val, ok = c.Get("bbb")
		require.True(t, ok)
		require.Equal(t, "filename2", val.FileName)

		wasInCache = c.Set("aaa", createItemStub("aaa", "filename3", 10))
		require.True(t, wasInCache)

		val, ok = c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, "filename3", val.FileName)

		val, ok = c.Get("ccc")
		require.False(t, ok)
		require.Nil(t, val)
	})

	t.Run("purge logic - sequentially", func(t *testing.T) {
		c := CreateCacheLRU(20, func(item *Item) {})

		c.Set("aaa", createItemStub("aaa", "filename1", 10))
		c.Set("bbb", createItemStub("bbb", "filename2", 10))
		c.Set("ccc", createItemStub("ccc", "filename3", 10))

		val, ok := c.Get("aaa")
		require.False(t, ok)
		require.Nil(t, val)

		val, ok = c.Get("bbb")
		require.True(t, ok)
		require.Equal(t, "filename2", val.FileName)

		val, ok = c.Get("ccc")
		require.True(t, ok)
		require.Equal(t, "filename3", val.FileName)
	})

	t.Run("purge logic - not useful", func(t *testing.T) {
		c := CreateCacheLRU(20, func(item *Item) {})

		c.Set("aaa", createItemStub("aaa", "filename1", 10))
		c.Set("bbb", createItemStub("bbb", "filename2", 10))
		c.Set("ccc", createItemStub("ccc", "filename3", 10))

		c.Get("bbb")
		c.Set("ddd", createItemStub("ddd", "filename4", 10))

		val, ok := c.Get("aaa")
		require.False(t, ok)
		require.Nil(t, val)

		val, ok = c.Get("ccc")
		require.False(t, ok)
		require.Nil(t, val)

		val, ok = c.Get("bbb")
		require.True(t, ok)
		require.Equal(t, "filename2", val.FileName)

		val, ok = c.Get("ddd")
		require.True(t, ok)
		require.Equal(t, "filename4", val.FileName)
	})

	t.Run("purge logic - many lightweight elements", func(t *testing.T) {
		c := CreateCacheLRU(20, func(item *Item) {})

		c.Set("aaa", createItemStub("aaa", "filename1", 5))
		c.Set("bbb", createItemStub("bbb", "filename2", 5))
		c.Set("ccc", createItemStub("ccc", "filename3", 5))
		c.Set("ddd", createItemStub("ddd", "filename4", 5))

		val, ok := c.Get("ccc")
		require.True(t, ok)
		require.Equal(t, "filename3", val.FileName)

		c.Set("eee", createItemStub("eee", "filename5", 15))
		val, ok = c.Get("eee")
		require.True(t, ok)
		require.Equal(t, "filename5", val.FileName)

		val, ok = c.Get("aaa")
		require.False(t, ok)
		require.Nil(t, val)

		val, ok = c.Get("bbb")
		require.False(t, ok)
		require.Nil(t, val)
	})

	t.Run("purge logic error - element has oversize", func(t *testing.T) {
		c := CreateCacheLRU(20, func(item *Item) {})
		c.Set("aaa", createItemStub("aaa", "filename1", 25))

		val, ok := c.Get("aaa")
		require.False(t, ok)
		require.Nil(t, val)
	})
}

func TestCacheMultithreading(_ *testing.T) {
	c := CreateCacheLRU(1000, func(item *Item) {})
	wg := &sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000_000; i++ {
			cacheKey := strconv.Itoa(i)
			c.Set(cacheKey, createItemStub(cacheKey, "filename"+cacheKey, 10))
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000_000; i++ {
			cacheKey := strconv.Itoa(rand.Intn(1_000_000))
			c.Set(cacheKey, createItemStub(cacheKey, "filename"+cacheKey, 10))
		}
	}()

	wg.Wait()
}

func createItemStub(key string, filename string, size uint64) *Item {
	return &Item{
		key:      key,
		FileName: filename,
		Size:     size,
	}
}
