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
)

type LocalCache interface {
	Get(key string) ([]byte, error)
	Set(key string, data []byte, expire time.Duration) error
	noLock
}

type noLock interface {
	NLGet(key string) ([]byte, error)
	NLSet(key string, data []byte, expire time.Duration) error
}
