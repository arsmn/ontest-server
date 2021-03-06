package generate

import (
	"crypto/rand"
	"math/big"
)

var rander = rand.Reader // random function

var (
	AlphaNum      = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	Alpha         = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	AlphaLowerNum = []rune("abcdefghijklmnopqrstuvwxyz0123456789")
	AlphaUpperNum = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	AlphaLower    = []rune("abcdefghijklmnopqrstuvwxyz")
	AlphaUpper    = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ")
	Numeric       = []rune("0123456789")
)

func RuneSequence(l int, allowedRunes []rune) (seq []rune, err error) {
	c := big.NewInt(int64(len(allowedRunes)))
	seq = make([]rune, l)

	for i := 0; i < l; i++ {
		r, err := rand.Int(rander, c)
		if err != nil {
			return seq, err
		}
		rn := allowedRunes[r.Uint64()]
		seq[i] = rn
	}

	return seq, nil
}

func RandomString(l int, allowedRunes []rune) string {
	seq, err := RuneSequence(l, allowedRunes)
	if err != nil {
		panic(err)
	}
	return string(seq)
}

func RandomInt(limit int64) (int64, error) {
	int, err := rand.Int(rand.Reader, big.NewInt(limit))
	if err != nil {
		return 0, err
	}
	return int.Int64(), nil
}

func ResetPasswordCode(email string) string {
	return RandomString(30, AlphaNum) + EncodeMD5(email)
}

func VerificationCode(email string) string {
	return RandomString(30, AlphaNum) + EncodeMD5(email)
}

func RandCode() string {
	return RandomString(10, AlphaNum)
}

func SessionToken() string {
	return RandomString(35, AlphaNum)
}

func VerifyResetPasswordCode(code, email string) bool {
	if len(code) <= 30 {
		return false
	}

	if EncodeMD5(email) != code[30:] {
		return false
	}

	return true
}

func TruncateString(str string, limit int) string {
	if len(str) < limit {
		return str
	}
	return str[:limit]
}
