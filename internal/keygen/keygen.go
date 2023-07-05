package keygen

import (
	"math/rand"
	"time"
)

func Generate() string {
	limit := len(words) - 1
	output := ""
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 4; i++ {
		randomNumber := rng.Intn(limit)
		output += words[randomNumber]
		if i != 3 {
			output += "-"
		}
	}
	if len(output) >= 16 {
		return output
	}

	return Generate()
}
