package concurrentmap

import (
	"math"
	"sync/atomic"
)

//the interface of map for concurrency security
type ConcurrentMap interface {
	//return the number of concurrency
	Concurrency() int
	//the value of element cannot be nil
	Put(key string, element interface{}) (bool, error)
	//return nil if the value not exist
	Get(key string) interface{}

	Delete(key string) bool

	Len() uint64
}

type PairRedistributor interface {
}

type myConcurrentMap struct {
	concurrency int
	segments    []Segment
	total       uint64
}

const (
	MAX_CONCURRENCY       = 16
	DEFAULT_BUCKET_NUMBER = 16
)

func NewConcurrentMap(concurrency int, pairRedistributor PairRedistributor) (ConcurrentMap, error) {
	if concurrency <= 0 {
		return nil, newIllegalParameterError("concurrency is to small")
	}
	if concurrency > MAX_CONCURRENCY {
		return nil, newIllegalParameterError("concurrency is to large")
	}
	cmap := &myConcurrentMap{}
	cmap.concurrency = concurrency
	cmap.segments = make([]Segment, concurrency)
	for i := 0; i < concurrency; i++ {
		cmap.segments[i] = newSegment(DEFAULT_BUCKET_NUMBER, pairRedistributor)
	}
	return cmap, nil
}

func newSegment(size int, redistributor PairRedistributor) Segment {

}

func (cmap *myConcurrentMap) Put(key string, element interface{}) (bool, error) {
	p, err := newPair(key, element)
	if err != nil {
		return false, err
	}
	s := cmap.findSegment(p.Hash())
	ok, err := s.Put(p)
	if ok {
		atomic.AddUint64(&cmap.total, 1)
	}
	return ok, err
}

func (cmap *myConcurrentMap) Get(key string) interface{} {
	keyHash := hash(key)
	s := cmap.findSegment(keyHash)
	pair := s.GetWithHash(key, keyHash)
	if pair == nil {
		return nil
	}
	return pair.Element()
}

func (cmap *myConcurrentMap) Delete(key string) bool {
	s := cmap.findSegment(hash(key))
	if s.Delete(key) {
		atomic.AddUint64(&cmap.total, ^uint64(0))
		return true
	}
	return false
}

//segment location algorithm
func (cmap *myConcurrentMap) findSegment(keyHash uint64) Segment {
	if cmap.concurrency == 1 {
		return cmap.segments[0]
	}
	var keyHashHigh int
	if keyHash > math.MaxUint32 {
		keyHashHigh = int(keyHash >> 48)
	} else {
		keyHashHigh = int(keyHash >> 16)
	}
	return cmap.segments[keyHashHigh%cmap.concurrency]
}
