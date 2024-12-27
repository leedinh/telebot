package bloomfilter

import (
	"hash"
	"hash/fnv"
)

type BloomFilter struct {
	bitArray     []bool
	size         int
	hashFuncions []hash.Hash
}

func NewBloomFilter(size int, numHashes int) *BloomFilter {
	bf := &BloomFilter{
		bitArray:     make([]bool, size),
		size:         size,
		hashFuncions: make([]hash.Hash, numHashes),
	}

	for i := 0; i < numHashes; i++ {
		bf.hashFuncions[i] = fnv.New64a()
	}

	return bf
}

func (bf *BloomFilter) Add(data []byte) {
	for _, hashfn := range bf.hashFuncions {
		hashfn.Reset()
		hashfn.Write(data)
		hash := hashfn.Sum(nil)
		index := int(hash[0]) % bf.size
		bf.bitArray[index] = true
	}
}

func (bf *BloomFilter) Contains(data []byte) bool {
	for _, hashfn := range bf.hashFuncions {
		hashfn.Reset()
		hashfn.Write(data)
		hash := hashfn.Sum(nil)
		index := int(hash[0]) % bf.size
		if !bf.bitArray[index] {
			return false
		}
	}

	return true
}
