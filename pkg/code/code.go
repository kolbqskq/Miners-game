package code

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

func Generate() string {
	code, _ := rand.Int(rand.Reader, big.NewInt(10000))
	return fmt.Sprintf("%04d", code.Int64())
}