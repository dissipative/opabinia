package cache

import (
	"testing"
)

func TestCache_Set(t *testing.T) {
	tests := []struct {
		name               string
		cache              func() *Cache
		expectedStorageLen int
	}{
		{
			name: "SetToZeroCache",
			cache: func() *Cache {
				return NewCache(0)
			},
			expectedStorageLen: 0,
		},
		{
			name: "SetToMinimumCache",
			cache: func() *Cache {
				return NewCache(1)
			},
			expectedStorageLen: 1,
		},
		{
			name: "SetToEmptyCache",
			cache: func() *Cache {
				return NewCache(5)
			},
			expectedStorageLen: 1,
		},
		{
			name: "SetToNonEmptyCacheWithNewKey",
			cache: func() *Cache {
				c := NewCache(5)
				c.Set("key1", "val1")
				return c
			},
			expectedStorageLen: 2,
		},
		{
			name: "SetToNonEmptyCacheWithExistingKey",
			cache: func() *Cache {
				c := NewCache(5)
				c.Set("key1", "val1")
				c.Set("testKey", "val2")
				return c
			},
			expectedStorageLen: 2,
		},
		{
			name: "SetToFullCache",
			cache: func() *Cache {
				c := NewCache(2)
				c.Set("key1", "val1")
				c.Set("key2", "val2")
				c.Set("key3", "val3")
				return c
			},
			expectedStorageLen: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := tt.cache()
			c.Set("testKey", "testVal")

			if len(c.storage) != tt.expectedStorageLen {
				t.Errorf("Expected storage length is %d, got %d", tt.expectedStorageLen, len(c.storage))
			}

			if tt.expectedStorageLen == 0 {
				return
			}

			found, ok := c.storage["testKey"]
			if !ok {
				t.Error("Expected `testKey` to be found")
			}

			v, ok := found.Value.(*item).val.(string)
			if !ok || v != "testVal" {
				t.Error("Expected value to be `testVal`")
			}
		})
	}
}
