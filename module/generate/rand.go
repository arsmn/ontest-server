package generate

import (
	"crypto/rand"
	"math/big"
)

const (
	numbers  = "0123456789"
	alphanum = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
)

func RandomInt(max *big.Int) (int, error) {
	rand, err := rand.Int(rand.Reader, max)
	if err != nil {
		return 0, err
	}

	return int(rand.Int64()), nil
}

func RandomChars(n int, src string) (string, error) {
	buffer := make([]byte, n)
	max := big.NewInt(int64(len(src)))

	for i := 0; i < n; i++ {
		index, err := RandomInt(max)
		if err != nil {
			return "", err
		}

		buffer[i] = src[index]
	}

	return string(buffer), nil
}

func RandomString(n int) (string, error) {
	return RandomChars(n, alphanum)
}

func RandomNumber(n int) (string, error) {
	return RandomChars(n, numbers)
}
