package main

import (
	"fmt"
	// murmur3 使用这个计算 hash，嘎嘎快
	"github.com/spaolacci/murmur3"
	"hash"
	"hash/fnv"
)

type HashFunc func() hash.Hash64

// DefaultHash 默认提供 4个 hash 函数
var DefaultHash = []HashFunc{
	func() hash.Hash64 { return fnv.New64() },
	func() hash.Hash64 { return murmur3.New64() },
	func() hash.Hash64 { return NewMD5() },
	func() hash.Hash64 {
		return NewSha256()
	},
}

type BloomFilter struct {
	keys   []byte     // 存储布隆过滤器的数组
	k      uint64     // hash 函数的个数
	size   uint64     // 数组的大小
	hashes []HashFunc // hash 函数的数组
}

func NewBloomFilter(k, size uint64, hashes []HashFunc) (*BloomFilter, error) {
	// 如果没有提供 hash 函数，就使用默认实现
	if len(hashes) == 0 {
		hashes = DefaultHash
	}

	// 如果指定的 hash 数量大于已有的 hash 函数，干它，直接报错，不惯着
	if k > uint64(len(hashes)) {
		return nil, fmt.Errorf("the value of K must be smaller than the number of the Hash function")
	}
	return &BloomFilter{
		keys:   make([]byte, size),
		k:      k,
		size:   size,
		hashes: DefaultHash,
	}, nil
}

// computeSlot 通过 hash 函数计算数组下标
func (b BloomFilter) computeSlot(index uint64, data []byte) (uint64, error) {
	hashFunc := b.hashes[index]()
	hashFunc.Reset()
	_, err := hashFunc.Write(data)
	if err != nil {
		return 0, err
	}
	// 取余取余取余，重点要考
	slot := hashFunc.Sum64() % b.size
	return slot, nil

}

// Add 添加一个 元素
func (b BloomFilter) Add(data []byte) error {
	for i := uint64(0); i <= b.k; i++ {
		slot, err := b.computeSlot(i, data)
		if err != nil {
			return err
		}
		// 将计算的位置设置为 1
		b.keys[slot] = 1
	}
	return nil
}

// IsExists 判断是否存在
func (b BloomFilter) IsExists(data []byte) (bool, error) {
	for i := uint64(0); i <= b.k; i++ {
		slot, err := b.computeSlot(i, data)
		if err != nil {
			return false, err
		}
		// 出现一个 0, 说明一定不存在
		if b.keys[slot] == 0 {
			return false, nil
		}
	}
	return true, nil
}
