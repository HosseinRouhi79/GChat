package helper

import (
	"math/rand"
	"strconv"
	"time"
)
func GeneratePassword() string {
	source := rand.NewSource(time.Now().Unix())
	rng := rand.New(source)

	max := 9999
	min := 1000
	randNumber := rng.Intn((max - min)) + min

	return strconv.Itoa(randNumber)
}