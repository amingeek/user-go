package cache_test

import (
	"testing"
	"time"

	"user-go/internal/cache"

	"github.com/stretchr/testify/assert"
)

func TestInMemoryCache_SetGetDelete(t *testing.T) {
	c := cache.NewInMemoryCache()

	// تست SetWithTTL و Get
	err := c.SetWithTTL("key1", "value1", 1) // 1 ثانیه TTL
	assert.NoError(t, err)

	val, err := c.Get("key1")
	assert.NoError(t, err)
	assert.Equal(t, "value1", val)

	// تست Delete
	err = c.Delete("key1")
	assert.NoError(t, err)

	_, err = c.Get("key1")
	assert.Error(t, err)
}

func TestInMemoryCache_Expire(t *testing.T) {
	c := cache.NewInMemoryCache()

	err := c.SetWithTTL("key2", "val2", 1) // یک ثانیه TTL
	assert.NoError(t, err)

	time.Sleep(1100 * time.Millisecond) // کمی بیشتر از 1 ثانیه صبر می‌کنیم

	_, err = c.Get("key2")
	assert.Error(t, err) // باید expired باشه و ارور بده
}

func TestInMemoryCache_IncrWithExpire_NewKey(t *testing.T) {
	c := cache.NewInMemoryCache()

	val, err := c.IncrWithExpire("counter", 1)
	assert.NoError(t, err)
	assert.Equal(t, 1, val)
}

func TestInMemoryCache_IncrWithExpire_ExistingKey(t *testing.T) {
	c := cache.NewInMemoryCache()

	val, err := c.IncrWithExpire("counter", 5)
	assert.NoError(t, err)
	assert.Equal(t, 1, val)

	val, err = c.IncrWithExpire("counter", 5)
	assert.NoError(t, err)
	assert.Equal(t, 2, val)

	// تست expire: مقدار نباید ریست شود چون هنوز TTL منقضی نشده
	val, err = c.IncrWithExpire("counter", 5)
	assert.NoError(t, err)
	assert.Equal(t, 3, val)
}

func TestInMemoryCache_IncrWithExpire_Expired(t *testing.T) {
	c := cache.NewInMemoryCache()

	val, err := c.IncrWithExpire("counter_exp", 1)
	assert.NoError(t, err)
	assert.Equal(t, 1, val)

	time.Sleep(1100 * time.Millisecond)

	val, err = c.IncrWithExpire("counter_exp", 1)
	assert.NoError(t, err)
	assert.Equal(t, 1, val) // چون expired شده دوباره باید از اول بشماریم
}
