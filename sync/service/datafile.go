package service

import (
	"errors"
	"io"
	"os"
	"sync"
	"sync/atomic"
)

//DataFile indicates the interface config  of data file
//multiple read operations should be in order and cannot be repeated;
//every data block is in the same size,if the data to be written exceeds this value,
//the excess will be truncation
type DataFile interface {
	//read a data block
	Read() (rsn int64, d Data, err error)
	//write a data block
	Write(d Data) (wsn int64, err error)
	//get the last read data block's serial number
	RSN() int64

	WSN() int64

	DataLen() uint32

	Close() error
}

type Data []byte

var pool sync.Pool

type myDataFile struct {
	f *os.File
	//read and write lock for operating file
	fmutex sync.RWMutex
	//indicates write progress
	woffset int64
	roffset int64
	//lock for update woffset
	wmutex sync.Mutex
	//lock for update roffset
	rmutex sync.Mutex
	//the size of data block
	dataLen uint32

	rcond *sync.Cond
}

func NewDataFile(path string, datalen uint32) (DataFile, error) {
	f, err := os.Create(path)
	if err != nil {
		return nil, err
	}
	if datalen == 0 {
		return nil, errors.New("Invalid data length!")
	}
	df := &myDataFile{f: f, dataLen: datalen}
	df.rcond = sync.NewCond(df.fmutex.RLocker())
	return df, nil
}

func (df *myDataFile) Read() (rsn int64, d Data, err error) {
	//read and update roffset
	var offset int64
	/*	df.rmutex.Lock()
		offset = df.roffset
		df.roffset += int64(df.dataLen)
		df.rmutex.Unlock()*/
	for {
		offset = atomic.LoadInt64(&df.roffset)
		if atomic.CompareAndSwapInt64(&df.roffset, offset, (offset + int64(df.dataLen))) {
			break
		}
	}

	//read a data block
	rsn = offset / int64(df.dataLen)
	bytes := make([]byte, df.dataLen)

	/*
		df.fmutex.RLock()
		defer df.fmutex.RUnlock()
		_,err  = df.f.ReadAt(bytes,offset)
		if err != nil{
			return
		}
		d = bytes*/

	/*
		for {
			df.fmutex.RLock()
			_,err = df.f.ReadAt(bytes,offset)
			if err != nil {
				//dealing with boundary conditions
				if err == io.EOF {
					df.fmutex.RUnlock()
					continue
				}
				df.fmutex.RUnlock()
				return
			}
			d = bytes
			df.fmutex.RUnlock()
			return
		}
		return*/

	df.fmutex.RLock()
	defer df.fmutex.RUnlock()
	for {
		_, err = df.f.ReadAt(bytes, offset)
		if err != nil {
			if err == io.EOF {
				df.rcond.Wait()
				continue
			}
			return
		}
		d = bytes
		return
	}

}

func (df *myDataFile) Write(d Data) (wsn int64, err error) {
	//read and update woffset
	df.wmutex.Lock()
	offset := df.woffset
	df.woffset += int64(df.dataLen)
	df.wmutex.Unlock()

	//write a data block
	wsn = offset / int64(df.dataLen)
	var bytes []byte
	if len(d) > int(df.dataLen) {
		bytes = d[0:df.dataLen]
	} else {
		bytes = d
	}
	df.fmutex.Lock()
	defer df.fmutex.Unlock()
	_, err = df.f.Write(bytes)
	df.rcond.Signal()
	return
}

func (df *myDataFile) RSN() int64 {
	/*	df.rmutex.Lock()
		defer df.rmutex.Unlock()*/
	offset := atomic.LoadInt64(&df.roffset)
	return offset / int64(df.dataLen)
}

func (df *myDataFile) WSN() int64 {
	df.wmutex.Lock()
	defer df.wmutex.Unlock()
	return df.woffset / int64(df.dataLen)
}

func (df *myDataFile) DataLen() uint32 {
	return df.dataLen
}

func (df *myDataFile) Close() error {
	return df.f.Close()
}
