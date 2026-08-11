package utils

import (
	"math/rand"
	"fmt"
)

func GenerateOTP(length int) (int, error) {
	if length <= 0 {
		return 0, fmt.Errorf("invalid OTP length")
	}
	otp := 0
	for i := 0; i < length; i++ {
		otp = otp*10 + rand.Intn(10)
	}
	return otp, nil
}
