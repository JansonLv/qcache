package cache

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/coocood/freecache"
	"github.com/stretchr/testify/assert"
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

func Test_cacheRepository_GetOrSetDataFunc(t *testing.T) {
	repo := NewCacheRepository()
	id := 0
	key := fmt.Sprintf(userKey, id)
	ioInfo := getInfoByIO(id, "jansonlv")
	func() {
		// 先获取缓存，缓存不存在，设置缓存
		info := &UserInfo{}
		setDataFunc, err := repo.GetOrSetDataFunc(key, info)
		// 缓存为找到
		assert.Equal(t, err, freecache.ErrNotFound)
		defer func() {
			// 设置缓存
			_ = setDataFunc(info)
		}()
		info = ioInfo
	}()
	func() {
		info := &UserInfo{}
		_, err := repo.GetOrSetDataFunc(key, info)
		// 缓存获取成功
		assert.NoError(t, err)
		// 判断是否是ioInfo中的缓存还是空数据
		assert.Equal(t, info, ioInfo)
	}()
	func() {
		info := &UserInfo{}
		time.Sleep(time.Second * 3)
		// 三秒钟后重新获取缓存
		_, err := repo.GetOrSetDataFunc(key, info)
		assert.Equal(t, err, freecache.ErrNotFound)
	}()
}

func Test_cacheRepository_GetOrSetDataFuncWithConditionOption(t *testing.T) {
	repo := NewCacheRepository()
	id := 20
	key := fmt.Sprintf(userKey, id)
	ioInfo := getInfoByIO(id, "jansonlv")
	func() {
		info := &UserInfo{}
		setDataFunc, err := repo.GetOrSetDataFunc(key, info, WithConditionOption(id < 10))
		// 缓存不满足条件
		assert.Equal(t, err, ConditionErr)
		defer func() {
			// 设置缓存
			_ = setDataFunc(info)
		}()
		info = ioInfo
	}()
	func() {
		info := &UserInfo{}
		_, err := repo.GetOrSetDataFunc(key, info)
		// 缓存条件不满足，未获取到缓存
		assert.Equal(t, err, freecache.ErrNotFound)
	}()
}

func Test_cacheRepository_GetOrSetDataWithExpireOption(t *testing.T) {
	repo := NewCacheRepository()
	id := 0
	key := fmt.Sprintf(userKey, id)
	ioInfo := getInfoByIO(id, "jansonlv")
	func() {
		// 先获取缓存，缓存不存在，设置缓存
		info := &UserInfo{}
		setDataFunc, err := repo.GetOrSetDataFunc(key, info, WithExpireOption(time.Second*5))
		// 缓存为找到
		assert.Equal(t, err, freecache.ErrNotFound)
		defer func() {
			// 设置缓存
			_ = setDataFunc(info)
		}()
		info = ioInfo
	}()
	func() {
		info := &UserInfo{}
		_, err := repo.GetOrSetDataFunc(key, info)
		// 缓存获取成功
		assert.NoError(t, err)
		// 判断是否是ioInfo中的缓存还是空数据
		assert.Equal(t, info, ioInfo)
	}()
	func() {
		info := &UserInfo{}
		time.Sleep(time.Second * 3)
		// 三秒钟后重新获取缓存
		_, err := repo.GetOrSetDataFunc(key, info)
		// 缓存获取成功
		assert.NoError(t, err)
		// 判断是否是ioInfo中的缓存还是空数据
		assert.Equal(t, info, ioInfo)
	}()
	func() {
		info := &UserInfo{}
		time.Sleep(time.Second * 2)
		// 再2秒钟后重新获取缓存，此时缓存失效
		_, err := repo.GetOrSetDataFunc(key, info)
		assert.Equal(t, err, freecache.ErrNotFound)
	}()
}

func Test_cacheRepository_GetOrSetDataWithFunc(t *testing.T) {
	repo := NewCacheRepository()
	id := 0
	key := fmt.Sprintf(userKey, id)
	ioInfo := getInfoByIO(id, "jansonlv")

	func() {
		// 先获取缓存，缓存不存在，设置缓存
		info := &UserInfo{}
		err := repo.GetWithMarshal(
			key,
			info,
			func() (interface{}, error) {
				return getInfoByIO(id, "jansonlv"), nil
			},
		)
		assert.NoError(t, err)
		assert.Equal(t, info, ioInfo)
	}()
	func() {
		info := &UserInfo{}
		err := repo.GetWithMarshal(
			key,
			info,
			func() (interface{}, error) {
				return getInfoByIO(id, "jansonlv1"), nil
			},
		)
		assert.NoError(t, err)
		assert.Equal(t, info, ioInfo)
	}()
	func() {
		id := 1
		var err1 = errors.New("aa")
		info := &UserInfo{}
		key := fmt.Sprintf(userKey, id)
		err := repo.GetWithMarshal(
			key,
			info,
			func() (interface{}, error) {
				return getInfoByIO(id, "jansonlv111"), err1
			},
		)
		assert.Equal(t, err, err1)
		assert.Equal(t, info, &UserInfo{})
	}()
	func() {
		time.Sleep(time.Second * 3)
		info := &UserInfo{}
		err := repo.GetWithMarshal(
			key,
			info,
			func() (interface{}, error) {
				return getInfoByIO(id, "jansonlv2"), nil
			},
		)
		assert.NoError(t, err)
		assert.NotEqual(t, info, ioInfo)
		assert.Equal(t, info, getInfoByIO(id, "jansonlv2"))
	}()
	//func() {
	//	info := &UserInfo{}
	//	time.Sleep(time.Second * 2)
	//	// 再2秒钟后重新获取缓存，此时缓存失效
	//	_, err := repo.GetOrSetDataFunc(key, info)
	//	assert.Equal(t, err, freecache.ErrNotFound)
	//}()
}

func Test_cacheRepository_GetOrSetDataWithFunc1(t *testing.T) {
	repo := NewCacheRepository()
	id := 0
	key := fmt.Sprintf(userKey, id)
	ioInfo := getInfoByIO(id, "jansonlv")

	func() {
		// 先获取缓存，缓存不存在，设置缓存
		info := &UserInfo{}
		err := repo.GetWithCopier(
			key,
			info,
			func() (interface{}, error) {
				return getInfoByIO2(id, "jansonlv"), nil
			},
		)
		assert.NoError(t, err)
		assert.Equal(t, info, ioInfo)
	}()
	func() {
		info := &UserInfo{}
		err := repo.GetWithCopier(
			key,
			info,
			func() (interface{}, error) {
				return getInfoByIO(id, "jansonlv1"), nil
			},
		)
		assert.NoError(t, err)
		assert.Equal(t, info, ioInfo)
	}()
	func() {
		id := 1
		var err1 = errors.New("aa")
		info := &UserInfo{}
		key := fmt.Sprintf(userKey, id)
		err := repo.GetWithCopier(
			key,
			info,
			func() (interface{}, error) {
				return getInfoByIO(id, "jansonlv111"), err1
			},
		)
		assert.Equal(t, err, err1)
		assert.Equal(t, info, &UserInfo{})
	}()
	func() {
		id := 10
		info := &UserInfo{}
		key := fmt.Sprintf(userKey, id)
		err := repo.GetWithCopier(
			key,
			info,
			func() (interface{}, error) {
				return getInfoByIO(id, "jansonlv10-1"), nil
			},
			WithConditionOption(id <= 5),
		)
		assert.NoError(t, err)
		assert.Equal(t, info, getInfoByIO(id, "jansonlv10-1"))
	}()
	func() {
		// 该id和上面的id一样，缓存key是一样的，但是不会缓存，因此，都会使用同一个key处理
		id := 10
		info := &UserInfo{}
		key := fmt.Sprintf(userKey, id)
		err := repo.GetWithCopier(
			key,
			info,
			func() (interface{}, error) {
				return getInfoByIO(id, "jansonlv10-2"), nil
			},
			WithConditionOption(id <= 5),
		)
		assert.NoError(t, err)
		assert.Equal(t, info, getInfoByIO(id, "jansonlv10-2"))
	}()
	func() {
		time.Sleep(time.Second * 3)
		info := &UserInfo{}
		err := repo.GetWithCopier(
			key,
			info,
			func() (interface{}, error) {
				return getInfoByIO(id, "jansonlv2"), nil
			},
		)
		assert.NoError(t, err)
		assert.NotEqual(t, info, ioInfo)
		assert.Equal(t, info, getInfoByIO(id, "jansonlv2"))
	}()
	//func() {
	//	info := &UserInfo{}
	//	time.Sleep(time.Second * 2)
	//	// 再2秒钟后重新获取缓存，此时缓存失效
	//	_, err := repo.GetOrSetDataFunc(key, info)
	//	assert.Equal(t, err, freecache.ErrNotFound)
	//}()
}


func Test_cacheRepository_Builder(t *testing.T) {
	repo := NewCacheRepository()
	id := 1
	key := fmt.Sprintf(userKey, id)
	ioInfo := getInfoByIO(id, "jansonlv")
	value := UserInfo{}
	err := repo.GetCache(key).Find(&value)
	assert.Error(t, err)
	err = repo.GetCache(key).SetGetDataFunc(func() (interface{}, error) {
		return getInfoByIO(id, "jansonlv"), nil
	}).Find(&value)
	assert.NoError(t, err)
	assert.Equal(t, &value, ioInfo)
	value2 := UserInfo{}
	err = repo.GetCache(key).Find(&value2)
	assert.NoError(t, err)
	assert.Equal(t, &value2, ioInfo)
}


func BenchmarkGetSet(b *testing.B) {
	repo1 := NewCacheRepository()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		id := i % 10
		ioInfo := &UserInfo{}
		key := fmt.Sprintf("name1:%d", id)
		if err := repo1.GetData(key, ioInfo); err == nil {
			continue
		}
		ioInfo = getInfoByIO(id, "jansonlv")
		_ = repo1.SetData(key, ioInfo)
	}
}

func BenchmarkGetOrSetFunc(b *testing.B) {
	repo2 := NewCacheRepository()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		func() {
			id := i % 10
			ioInfo := &UserInfo{}
			key := fmt.Sprintf("name2:%d", id)
			setFunc, err := repo2.GetOrSetDataFunc(key, ioInfo)
			if err == nil {
				return
			}
			defer func() {
				_ = setFunc(ioInfo)
			}()
			ioInfo = getInfoByIO(id, "jansonlv")
		}()
	}
}

func BenchmarkGetWithMarshal(b *testing.B) {
	repo2 := NewCacheRepository()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		id := i % 10
		ioInfo := &UserInfo{}
		key := fmt.Sprintf("name2:%d", id)
		_ = repo2.GetWithMarshal(key, ioInfo, func() (interface{}, error) {
			return getInfoByIO(id, "jansonlv"), nil
		})
	}
}

func BenchmarkGetWithCopier(b *testing.B) {
	repo2 := NewCacheRepository()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		id := i % 10
		ioInfo := &UserInfo{}
		key := fmt.Sprintf("name2:%d", id)
		_ = repo2.GetWithCopier(key, ioInfo, func() (interface{}, error) {
			return getInfoByIO(id, "jansonlv"), nil
		})
	}
}

/*
读写比例： 10：1
goos: darwin
goarch: amd64
pkg: github.com/JansonLv/go-cache
cpu: Intel(R) Core(TM) i7-1068NG7 CPU @ 2.30GHz
BenchmarkName1-8          993747              1203 ns/op             304 B/op          9 allocs/op
BenchmarkName2-8          916068              1302 ns/op             376 B/op         12 allocs/op
PASS
ok      github.com/JansonLv/go-cache    11.609s

*/

/*
读写比例：100：1
goos: darwin
goarch: amd64
pkg: github.com/JansonLv/go-cache
cpu: Intel(R) Core(TM) i7-1068NG7 CPU @ 2.30GHz
BenchmarkName1-8          965590              1229 ns/op             311 B/op          9 allocs/op
BenchmarkName2-8          895328              1339 ns/op             383 B/op         12 allocs/op
PASS
ok      github.com/JansonLv/go-cache    12.240s
*/
