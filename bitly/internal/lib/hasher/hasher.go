package hasher

import (
	"crypto/sha256"
	"math/big"
	"sync"

	"github.com/leedinh/telebot/bitly/internal/lib/bloomfilter"
)

// Base62 characters
const base62Chars = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

var mu sync.Mutex // Mutex for thread-safety

func toBase62(num *big.Int) string {
	var result string
	for num.Cmp(big.NewInt(0)) > 0 {
		remainder := new(big.Int)
		num.DivMod(num, big.NewInt(62), remainder)
		result = string(base62Chars[remainder.Int64()]) + result
	}

	return result
}

func toAlias(url string, size int) string {
	hash := sha256.New()
	hash.Write([]byte(url))
	hashedURL := hash.Sum(nil)

	hashInt := new(big.Int).SetBytes(hashedURL)
	alias := toBase62(hashInt)

	return alias[:size]
}

// Hasher is an interface that wraps the Hash method.

func GenerateAlias(size int, url string, bf *bloomfilter.BloomFilter) string {
	mu.Lock()

	defer mu.Unlock()

	for {
		alias := toAlias(url, size)
		if !bf.Contains([]byte(alias)) {
			return alias
		}

		url += "D"
	}

}
