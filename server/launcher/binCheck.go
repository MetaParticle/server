package main

import (
	"crypto/md5"
	"bytes"
)

func hash(b []byte) []byte {
	hasher := md5.New()
	hasher.Write(b)
	res := make([]byte, 0)
	return hasher.Sum(res)
}

func CompareHash2Bin(hashed []byte, bin []byte) bool {
	return bytes.Compare(hash(bin), hashed) == 0
}

func GetHashSize() int {
	return md5.Size
}

func compareHash(hash1, hash2 []byte) bool {
	return bytes.Compare(hash1, hash2) == 0
}
