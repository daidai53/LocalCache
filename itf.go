// Copyright@daidai53 2023
package cache

import (
	"errors"
	"time"
)

var (
	ErrCodeBadCache        = errors.New("cache is invalid")
	ErrCodeRecordNotFound  = errors.New("key is not found in cache")
	ErrCodeWrongBucketHash = errors.New("bucket hash is wrong")
	ErrCodeNoExpireTime    = errors.New("key has not expiration")
)

type LocalCache interface {
	Get(key string) ([]byte, error)
	Set(key string, data []byte, expire time.Duration) error
	noLock
}

type noLock interface {
	SafeOperate(key string, f func(c LocalCache) error) error
	NLTTL(key string) (int, error)
	NLGet(key string) ([]byte, error)
	NLSet(key string, data []byte, expire time.Duration) error
}
