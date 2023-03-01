package rand

import (
	"math/rand"
	"time"
)

func GetRand(max int) int64 {
	rand.Seed(time.Now().Unix())
	return rand.Int63n(int64(max))
}
