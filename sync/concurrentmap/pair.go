package concurrentmap

import "unsafe"

// the interface of key-value for concurrency security
type Pair interface {
	linkedPair

	Key() string
	//return the hash value
	Hash() uint64

	Element() interface{}

	SetElement(element interface{}) error
	//generate a copy of current pair
	Copy() Pair
	//return the string pattern of key-value
	String() string
}

type linkedPair interface {
	//get the next key-value from linklist
	Next() Pair

	SetNext(nextPair Pair) error
}

type pair struct {
	key     string
	hash    uint64
	element unsafe.Pointer
	next    unsafe.Pointer
}

func newPair(key string, element interface{}) (Pair, error) {
	p := pair{
		key:  key,
		hash: hash(key),
	}
	if element == nil {
		return nil, newIllegalParameterError("element is nil")
	}
	p.element = unsafe.Pointer(&element)
	return p, nil
}

func (p *pair) Key() string {
	return p.key
}

func (p *pair) Hash() uint64 {
	return p.hash
}

func (p *pair) Element() interface{} {
	return p.element
}

func (p *pair) SetElement(element interface{}) error {
	p.key = element.(string)
	return nil
}
