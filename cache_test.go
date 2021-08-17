package cache

import (
	"fmt"
	"github.com/JansonLv/go-cache/store"
	"github.com/coocood/freecache"
	"github.com/patrickmn/go-cache"
	"github.com/stretchr/testify/assert"
	"strconv"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

type UserInfo struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type UserInfo2 struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

func getInfoByIO(id int, name string) *UserInfo {
	return &UserInfo{
		Id:   id,
		Name: name,
	}
}

func getInfoByIO2(id int, name string) *UserInfo2 {
	return &UserInfo2{
		Id:   id,
		Name: name,
	}
}

const userKey = "userInfo:%d"

func Test_freecache(t *testing.T) {
	id := 1
	key := fmt.Sprintf(userKey, id)
	qCache := NewCacheRepository(store.NewFreeCache(freecache.NewCache(512 * 1024)))
	v1 := UserInfo{}
	err := qCache.GetCacheWithOptions(key, &v1, func() (interface{}, error) {
		return getInfoByIO(id, strconv.Itoa(id)), nil
	}, WithExpireOption(time.Hour))

	assert.NoError(t, err)
	assert.Equal(t, &v1, getInfoByIO(id, strconv.Itoa(id)))

	v2 := UserInfo{}
	err = qCache.GetKey(key).Find(&v2)
	assert.NoError(t, err)
	assert.Equal(t, &v1, &v2)

	id = 3
	key = fmt.Sprintf(userKey, id)
	err = qCache.GetKey(key).Find(&v1)
	assert.Error(t, err)

	err = qCache.GetKey(key).SetGetDataFunc(func() (interface{}, error) {
		return getInfoByIO(id, strconv.Itoa(id)), nil
	}).SetIsSave(id < 4).SetExpire(time.Second * 5).Find(&v1)
	assert.NoError(t, err)
	assert.Equal(t, &v1, getInfoByIO(id, strconv.Itoa(id)))

	err = qCache.GetKey(key).Find(&v1)
	assert.NoError(t, err)

	id = 4
	key = fmt.Sprintf(userKey, id)
	err = qCache.GetKey(key).SetGetDataFunc(func() (interface{}, error) {
		return getInfoByIO(id, strconv.Itoa(id)), nil
	}).SetIsSave(id < 4).Find(&v1)
	assert.NoError(t, err)
	assert.Equal(t, &v1, getInfoByIO(id, strconv.Itoa(id)))

	err = qCache.GetKey(key).Find(&v1)
	assert.Error(t, err)

	var int1 int
	key = "1"
	err = qCache.GetCacheWithOptions(key, &int1, func() (interface{}, error) {
		var b = 1
		return &b, nil
	})
	assert.NoError(t, err)
	assert.Equal(t, int1, 1)

	id = 2
	key = fmt.Sprintf(userKey, id)
	err = qCache.GetKey(key).Find(&v1)
	assert.Error(t, err)
	assert.NotEqual(t, &v1, getInfoByIO(id, strconv.Itoa(id)))
	id = 3
	key = fmt.Sprintf(userKey, id)
	err = qCache.GetKey(key).Find(&v1)
	assert.NoError(t, err)
	assert.Equal(t, &v1, getInfoByIO(id, strconv.Itoa(id)))
}

func Test_GoCache(t *testing.T) {
	id := 1
	key := fmt.Sprintf(userKey, id)
	client := store.NewGoCacheStore(cache.New(time.Minute*5, time.Second))
	qCache := NewCacheRepository(client)
	v1 := UserInfo{}
	err := qCache.GetCacheWithOptions(key, &v1, func() (interface{}, error) {
		return getInfoByIO(id, strconv.Itoa(id)), nil
	})
	assert.NoError(t, err)
	assert.Equal(t, &v1, getInfoByIO(id, strconv.Itoa(id)))
	v2 := UserInfo{}
	err = qCache.GetKey(key).Find(&v2)
	assert.NoError(t, err)
	assert.Equal(t, &v1, &v2)
	id = 3
	key = fmt.Sprintf(userKey, id)
	err = qCache.GetKey(key).Find(&v1)
	assert.Error(t, err)
	var int1 int
	key = "1"
	err = qCache.GetCacheWithOptions(key, &int1, func() (interface{}, error) {
		var b = 1
		return &b, nil
	})
	assert.NoError(t, err)
	assert.Equal(t, int1, 1)
}

// 请求归并处理
func Test_CacheSingleFlight(t *testing.T)  {
	id := 1
	key := fmt.Sprintf(userKey, id)
	qCache := NewCacheRepository(store.NewFreeCache(freecache.NewCache(512 * 1024)))
	v := getInfoByIO(id, strconv.Itoa(id))
	var wg sync.WaitGroup
	var count int64
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			v1 := UserInfo{}
			err := qCache.GetCacheWithOptions(key, &v1, func() (interface{}, error) {
				atomic.AddInt64(&count, 1)
				return getInfoByIO(id, strconv.Itoa(id)), nil
			}, WithExpireOption(time.Second))
			assert.NoError(t, err)
			assert.Equal(t, &v1, v)
		}()
	}
	wg.Wait()
	fmt.Println(count)

	time.Sleep(time.Second)
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			v1 := UserInfo{}
			err := qCache.GetCacheWithOptions(key, &v1, func() (interface{}, error) {
				atomic.AddInt64(&count, 1)
				return getInfoByIO(id, strconv.Itoa(id)), nil
			})
			assert.NoError(t, err)
			assert.Equal(t, &v1, v)
		}()
	}
	wg.Wait()
	fmt.Println(count)

}