// Copyright (c) 1992-1993 The Regents of the University of California.
// All rights reserved.

package utils

import "math/rand"

// RandomInit is used to initialize the pseudo-random generator with a deterministic seed
func RandomInit(seed int64) {
	rand.Seed(seed)
}

// Random returns a random integer
func Random() int {
	return rand.Int()
}
