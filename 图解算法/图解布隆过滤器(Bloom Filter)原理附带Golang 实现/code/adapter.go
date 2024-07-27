package main

import (
	"crypto/md5"
	"crypto/sha256"
	"hash"
)

func Bytes16ToUint64(b []byte) uint64 {
	return uint64(b[7]^b[15]) | uint64(b[6]^b[14])<<8 | uint64(b[5]^b[13])<<16 | uint64(b[4]^b[12])<<24 |
		uint64(b[3]^b[11])<<32 | uint64(b[2]^b[10])<<40 | uint64(b[1]^b[9])<<48 | uint64(b[0]^b[8])<<56
}

type MD5Adapter struct {
	hash.Hash
}

func (h MD5Adapter) Sum64() uint64 {
	b := h.Sum(nil)
	return Bytes16ToUint64(b)
}

func NewMD5() hash.Hash64 {
	m := &MD5Adapter{
		md5.New(),
	}
	return m
}

type Sha256Adapter struct {
	hash.Hash
}

func (h Sha256Adapter) Sum64() uint64 {
	b := h.Sum(nil)
	return Bytes16ToUint64(b)
}

func NewSha256() hash.Hash64 {
	m := &Sha256Adapter{
		sha256.New(),
	}
	return m
}
