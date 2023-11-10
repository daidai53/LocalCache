// Copyright@daidai53 2023
package cache

import "testing"

func Test_LocalCacheV1(t *testing.T) {
	_ = NewLocalCacheV1(256)
	a := 0
	for {
		a++
	}
}
