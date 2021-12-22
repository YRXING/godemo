package concurrentmap

import "sync"

type Bucket interface {
	//the lock parameter should not be passed into if it is locked before calling this function
	Put(p Pair, lock sync.Locker) (bool, error)
}
