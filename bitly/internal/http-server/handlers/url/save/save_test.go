package save_test

import (
	"fmt"
	"testing"

	"github.com/leedinh/telebot/bitly/internal/lib/bloomfilter"
	"github.com/leedinh/telebot/bitly/internal/lib/hasher"
)

func TestSaveURL(t *testing.T) {
	bf := bloomfilter.NewBloomFilter(1000, 3)
	alias := hasher.GenerateAlias(7, "https://www.google.com", bf)
	bf.Add([]byte(alias))
	fmt.Println(alias)
	alias1 := hasher.GenerateAlias(7, "https://www.google.com", bf)
	fmt.Println(alias1)

}
