package generators

import (
	"math/rand"
	"time"
)

func random(min, max int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max-min) + min
}

func PhoneNumber() uint {
	return uint(random(5550000000, 5559999999))
}
