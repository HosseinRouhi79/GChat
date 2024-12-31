package helper

import (
	"chat/config"
	"math"
	"math/rand"
	"time"
)

type Otp struct {
	Limiter    time.Duration
	ExpireTime time.Duration
	Digits     int
}

func GenerateOtp() int {
	cfg := config.GetConfig()
	// limiter := cfg.Otp.Limiter
	digits := cfg.Otp.Digits
	// expireTime := cfg.Otp.ExpireTime
	source := rand.NewSource(time.Now().UnixNano())
	rng := rand.New(source) //create independent random generator

	min := int(math.Pow(10, float64(digits-1)))
	max := int(math.Pow(10, float64(digits)) - 1)

	randomNumber := rng.Intn(max-min) + min //999999-100000 = 899999 ====> + min(100000) // Intn(n) give random num from zero to n

	return randomNumber

}
