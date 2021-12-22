package service

import (
	"errors"
	"sync/atomic"
)

//concurrent security array
type ConcurrentArray interface {
	Set(index int32, elem int) (err error)
	Get(index int32) (elem int, err error)
	Len() uint32
}

type concurrentArray struct {
	length uint32
	val    atomic.Value
}

func NewConcurrentArray(length uint32) ConcurrentArray {
	array := concurrentArray{}
	array.length = length
	array.val.Store(make([]int, array.length))
	return &array
}

func (array *concurrentArray) Set(index int32, elem int) (err error) {
	//if err = array.checkIndex(uint32(index));err!= nil{
	//	return
	//}
	//if err = array.checkValue(elem);err!= nil{
	//	return
	//}
	newArray := make([]int, array.length)
	copy(newArray, array.val.Load().([]int))
	newArray[index] = elem
	array.val.Store(newArray)
	return
}

func (array *concurrentArray) Get(index int32) (elem int, err error) {
	//if err = array.checkIndex(uint32(index));err!= nil{
	//	return
	//}
	//if err = array.checkValue(elem);err!= nil{
	//	return
	//}
	elem = array.val.Load().([]int)[index]
	return
}

func (array *concurrentArray) Len() uint32 {
	return array.length
}

func (array *concurrentArray) checkIndex(index uint32) error {
	if array.length <= index {
		return errors.New("invalid index")
	}
	return nil
}
