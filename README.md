# qcache 
### all in one

## plan

* [ ] 支持多个内存缓存库
  * [x] [go-cache](github.com/patrickmn/go-cache)
  * [x] [freecache](https://github.com/coocood/freecache)
  * [ ] [ristretto](https://github.com/dgraph-io/ristretto)
  * [ ] [bigcache](https://github.com/allegro/bigcache)
  * [ ] [fastcache](https://github.com/VictoriaMetrics/fastcache)
  * [ ] [goburrow](https://github.com/goburrow/cache)
  * [ ] [gcache](https://github.com/bluele/gcache)
  * [ ] [groupcache](https://github.com/golang/groupcache) (待调研)
* [x] 并发安全，防内存穿透（Singleflight）
* [ ] 支持redis数据存储 
* [ ] 支持protobuf数据序列化和反序列化
* [ ] 缓存命中率统计