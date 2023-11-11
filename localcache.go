// Copyright@daidai53 2023
package cache

import (
	"fmt"
	"hash/fnv"
	"sync"
	"time"
)

type (
	LocalCacheV1 struct {
		bucketDepth int
		bucket      []localCacheBucket
	}

	// 单个桶中存放的结构
	localCacheBucket struct {
		mu        sync.RWMutex
		bucketMap map[string]dataBlock
	}

	// 缓存的数据块
	dataBlock struct {
		key        string // 预留用于后续核查
		data       []byte
		mu         sync.RWMutex
		expireTime time.Time
	}
)

func NewLocalCacheV1(bucketDepth int) LocalCache {
	lc := &LocalCacheV1{
		bucketDepth: bucketDepth,
		bucket:      make([]localCacheBucket, bucketDepth),
	}
	for _, lcb := range lc.bucket {
		lcb.bucketMap = make(map[string]dataBlock)
	}
	go monitor(lc)
	return lc
}

func (l *LocalCacheV1) Get(key string) ([]byte, error) {
	hash := l.hash(key)
	if hash > len(l.bucket) {
		return []byte{}, ErrCodeWrongBucketHash
	}
	bucket := l.bucket[hash]
	if !bucket.valid() {
		return []byte{}, ErrCodeBadCache
	}
	bucket.mu.RLock()
	db, ok := bucket.bucketMap[key]
	bucket.mu.RUnlock()
	if !ok {
		return []byte{}, ErrCodeRecordNotFound
	}
	db.mu.RLock()
	defer db.mu.RUnlock()
	exp := db.expireTime
	now := time.Now()
	if exp.Unix() < now.Unix() {
		return []byte{}, ErrCodeRecordNotFound
	}
	return db.data, nil
}

func (l *LocalCacheV1) Set(key string, data []byte, expire time.Duration) error {
	hash := l.hash(key)
	if hash > len(l.bucket) {
		return ErrCodeWrongBucketHash
	}
	bucket := l.bucket[hash]
	if !bucket.valid() {
		return ErrCodeBadCache
	}
	bucket.mu.Lock()
	defer bucket.mu.Unlock()
	db, ok := bucket.bucketMap[key]
	if ok {
		db.data = data
		db.expireTime = time.Now().Add(expire)
	} else {
		bucket.bucketMap[key] = dataBlock{
			key:        key,
			data:       data,
			expireTime: time.Now().Add(expire),
		}
	}
	return nil
}

func (l *LocalCacheV1) SafeOperate(key string, f func(c LocalCache) error) error {
	hash := l.hash(key)
	if hash > len(l.bucket) {
		return ErrCodeWrongBucketHash
	}
	bucket := l.bucket[hash]
	if !bucket.valid() {
		return ErrCodeBadCache
	}
	bucket.mu.RLock()
	db, ok := bucket.bucketMap[key]
	bucket.mu.RUnlock()
	if !ok {
		return ErrCodeRecordNotFound
	}
	db.mu.Lock()
	defer db.mu.Unlock()
	return f(l)
}

func (l *LocalCacheV1) NLGet(key string) ([]byte, error) {
	hash := l.hash(key)
	if hash > len(l.bucket) {
		return []byte{}, ErrCodeWrongBucketHash
	}
	db, ok := l.bucket[hash].bucketMap[key]
	if !ok {
		return []byte{}, ErrCodeRecordNotFound
	}
	exp := db.expireTime
	now := time.Now()
	if exp.Unix() < now.Unix() {
		return []byte{}, ErrCodeRecordNotFound
	}
	return db.data, nil
}

func (l *LocalCacheV1) NLSet(key string, data []byte, expire time.Duration) error {
	hash := l.hash(key)
	if hash > len(l.bucket) {
		return ErrCodeWrongBucketHash
	}
	bucket := l.bucket[hash]
	if !bucket.valid() {
		return ErrCodeBadCache
	}
	db, ok := bucket.bucketMap[key]
	if ok {
		db.data = data
		db.expireTime = time.Now().Add(expire)
	} else {
		bucket.bucketMap[key] = dataBlock{
			key:        key,
			data:       data,
			expireTime: time.Now().Add(expire),
		}
	}
	return nil
}

func (l *LocalCacheV1) hash(key string) int {
	hash32 := fnv.New32()
	_, err := hash32.Write([]byte(key))
	if err != nil {
		return -1
	}
	return int(hash32.Sum32()) / l.bucketDepth
}

func (l *localCacheBucket) valid() bool {
	valid := false
	l.mu.RLock()
	if l.bucketMap != nil {
		valid = true
	}
	l.mu.RUnlock()
	return valid
}

func monitor(lc *LocalCacheV1) {
	ticker := time.NewTicker(1 * time.Second)
	for {
		<-ticker.C
		//todo:implement monitor
		fmt.Println("核查一次")
	}
}
