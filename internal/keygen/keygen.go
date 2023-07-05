package keygen

import (
	"math/rand"
	"time"
)

func Generate() string {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	output := generateNLengthWord(rng, 4) + "-" + generateNLengthWord(rng, 4) + "-" + generateNLengthWord(rng, 6)
	return output
}

func generateNLengthWord(rng *rand.Rand, n int) string {
	word := words[rng.Intn(len(words)-1)]
	if len(word) == n {
		return word
	}
	return generateNLengthWord(rng, n)
}
