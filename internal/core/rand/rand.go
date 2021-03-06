package rand

import (
	"crypto/rand"
	"fmt"
	"log"
	"math/big"
	"strings"
)

var max, _ = new(big.Int).SetString(strings.Repeat("9", 100), 10)

type Random = func() (*big.Int, error)

func CryptoSafeRandom() Random {
	return func() (*big.Int, error) {
		return rand.Int(rand.Reader, max)
	}
}

// ConstRand is useful for tests
func ConstRand(value *big.Int) Random {
	return func() (*big.Int, error) {
		return value, nil
	}
}

func FromToRandom(from *big.Int, to *big.Int, rand Random) Random {
	return func() (*big.Int, error) {
		if from.Cmp(to) >= 0 {
			return nil, fmt.Errorf("RANDOM: min=%v must be less than max=%v", from, to)
		}
		value, e := rand()
		if e != nil {
			return nil, e
		}
		diff := new(big.Int).Sub(to, from)
		shift := new(big.Int).Mod(value, diff)
		return new(big.Int).Add(shift, from), nil
	}
}

// ConditionalRandom returns first random value that matches the predicate
func ConditionalRandom(maxTry int, predicate func(*big.Int) bool, rand Random) Random {
	return func() (*big.Int, error) {
		for i := 0; i < maxTry; i++ {
			value, e := rand()
			if e != nil {
				return nil, e
			}
			if !predicate(value) {
				continue
			}
			return value, nil
		}
		return nil, fmt.Errorf("conditional random: max try count(%v) exceeded", maxTry)
	}
}

// CyclicRandom is useful for tests
func CyclicRandom(values ...*big.Int) Random {
	if len(values) == 0 {
		log.Fatalf("CyclicRandom: values cannot be empty")
	}

	i := 0
	size := len(values)
	return func() (*big.Int, error) {
		res := values[i]
		i = (i + 1) % size
		return res, nil
	}
}
