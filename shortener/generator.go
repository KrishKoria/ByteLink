package shortener

import (
	"crypto/sha256"
	"fmt"
	"github.com/itchyny/base58-go"
	"math/big"
	"os"
)

func sha256Hash(data string) []byte {
	algo := sha256.New()
	algo.Write([]byte(data))
	return algo.Sum(nil)
}

func base58Encoder(bytes []byte) string {
	encoding := base58.BitcoinEncoding
	encoded, err := encoding.Encode(bytes)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	return string(encoded)
}

func GenerateShortURL(initialLink string, userId string) string {
	urlHash := sha256Hash(initialLink + userId)
	num := new(big.Int).SetBytes(urlHash).Uint64()
	shortURL := base58Encoder([]byte(fmt.Sprintf("%d", num)))
	return shortURL[:8]
}
