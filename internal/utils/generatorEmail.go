package generatorEmail

import (
	"math/rand"
)

func GenerateOTP() int {
	const otpLength = 6
	otp := 0
	for i := 0; i < otpLength; i++ {
		otp = otp*10 + rand.Intn(10)
	}
	return otp
}
